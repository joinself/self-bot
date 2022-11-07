package components

import (
	"strings"

	selfsdk "github.com/joinself/self-go-sdk"
)

type WelcomeComponent struct {
	client                     *selfsdk.Client
	registeredUserWelcomeMsg   string
	unRegisteredUserWelcomeMsg string
}

type WelcomeComponentConfig struct {
	Client                     *selfsdk.Client
	RegisteredUserWelcomeMsg   string
	UnRegisteredUserWelcomeMsg string
}

func NewWelcomeComponent(cfg WelcomeComponentConfig) *WelcomeComponent {
	return &WelcomeComponent{
		client:                     cfg.Client,
		registeredUserWelcomeMsg:   cfg.RegisteredUserWelcomeMsg,
		unRegisteredUserWelcomeMsg: cfg.UnRegisteredUserWelcomeMsg,
	}
}

func (i *WelcomeComponent) RecordCommands(r CommandRecorder) {
	return
}

func (i *WelcomeComponent) AfterStartHook(r CommandRecorder) {
	i.client.ChatService().OnConnection(func(iss, status string) {
		parts := strings.Split(iss, ":")
		if len(parts) > 1 {
			iss = parts[0]
		}
		body := i.unRegisteredUserWelcomeMsg
		if r.IsUser(iss) {
			body = i.registeredUserWelcomeMsg
		}
		if body != "" {
			i.client.ChatService().Message([]string{iss}, body)
		}
	})
}
