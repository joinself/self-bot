package components

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/chat"
)

type WelcomeComponent struct {
	client                     *selfsdk.Client
	registeredUserWelcomeMsg   string
	unRegisteredUserWelcomeMsg string
	store                      Store
	componentName              string
	connections                []string
}

type WelcomeComponentConfig struct {
	Client                     *selfsdk.Client
	RegisteredUserWelcomeMsg   string
	UnRegisteredUserWelcomeMsg string
	Store                      Store
}

func NewWelcomeComponent(cfg WelcomeComponentConfig) *WelcomeComponent {
	wc := WelcomeComponent{
		client:                     cfg.Client,
		registeredUserWelcomeMsg:   cfg.RegisteredUserWelcomeMsg,
		unRegisteredUserWelcomeMsg: cfg.UnRegisteredUserWelcomeMsg,
		componentName:              "welcome_component",
	}
	wc.load()

	return &wc
}

func (i *WelcomeComponent) RecordCommands(r CommandRecorder) {
	i.recordGetConnections(r)
	return
}

func (i *WelcomeComponent) AfterStartHook(r CommandRecorder) {
	i.client.ChatService().OnConnection(func(iss, status string) {
		parts := strings.Split(iss, ":")
		if len(parts) > 1 {
			iss = parts[0]
		}

		// Save the user
		i.connections = append(i.connections, iss)
		i.save()

		body := i.unRegisteredUserWelcomeMsg
		if r.IsUser(iss) {
			body = i.registeredUserWelcomeMsg
		}
		if body != "" {
			i.client.ChatService().Message([]string{iss}, body)
		}
	})
}

func (i *WelcomeComponent) recordGetConnections(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "get_connections",
		Summary: "Returns a csv file with the list of connected users.",
		Callback: func(cmd string, cm *chat.Message) string {
			resp := "Find attached a list of objects"
			content := strings.Join(i.connections, "\n")

			objects := make([]chat.MessageObject, 0)
			objects = append(objects, chat.MessageObject{
				Name: "connections_new_" + time.Now().String() + ".csv",
				Data: []byte(content),
				Mime: "text/csv",
			})

			cm.Message(resp, chat.MessageOptions{
				Objects: objects,
			})

			return ""
		},
	})
}

func (i *WelcomeComponent) load() {
	if i.store == nil {
		return
	}

	content, err := i.store.Get(i.componentName)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(content, &i.connections)
	if err != nil {
		log.Println(err)
		return
	}
}

func (i *WelcomeComponent) save() {
	if i.store == nil {
		return
	}

	content, err := json.Marshal(i.connections)
	if err != nil {
		log.Println(err)
		return
	}
	i.store.Set(i.componentName, content)
}
