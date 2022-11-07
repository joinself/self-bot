package bot

import (
	"errors"
	"log"
	"strings"

	"github.com/joinself/self-bot/components"
	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/chat"
)

type Fallback func(*chat.Message) bool

type SelfBotConfig struct {
	Client            *selfsdk.Client
	MessagingFallback []Fallback
	Components        []components.CommandComponent
}

type SelfBot struct {
	client            *selfsdk.Client
	commands          map[string]components.CommandCallback
	help              map[string]string
	auth              components.AuthCallback
	messagingFallback []Fallback
	components        []components.CommandComponent
}

func NewSelfBot(cfg SelfBotConfig) *SelfBot {
	bot := SelfBot{
		client:   cfg.Client,
		commands: make(map[string]components.CommandCallback, 0),
		help:     make(map[string]string, 0),
		auth: func(string, string) bool {
			return true
		},
		messagingFallback: cfg.MessagingFallback,
		components:        cfg.Components,
	}

	for _, c := range bot.components {
		c.RecordCommands(&bot)
	}

	return &bot
}

func (s *SelfBot) RecordCommand(cmd components.Command) {
	s.commands[cmd.Name] = cmd.Callback
	s.help[cmd.Name] = cmd.Summary
}

func (s *SelfBot) Start() {
	s.client.ChatService().OnMessage(func(cm *chat.Message) {
		cmd, err := s.getCommandName(cm.Body)
		if err != nil {
			log.Println("could not get the command name")
			s.processNonCommandMessages(cm)
			return
		}

		command, ok := s.commands[cmd]
		if !ok {
			log.Println("command " + cmd + "not registered")
			s.processNonCommandMessages(cm)
			return
		}

		if !s.auth(cm.ISS, cmd) {
			log.Println("user doesn't have permissions to interact with this command")
			s.processNonCommandMessages(cm)
			return
		}

		resp := command(cmd, cm)
		if len(resp) == 0 {
			log.Println("command responded with an empty string, skipping response...")
			// s.processNonCommandMessages(cm)
			return
		}

		_, err = cm.Message(resp)
		if err != nil {
			log.Println("message failed")
			return
		}
	})

	for _, c := range s.components {
		c.AfterStartHook(s)
	}
}

func (s *SelfBot) GetCommandParams(body string) []string {
	return components.GetCommandParams(body)
}

func (s *SelfBot) SetAuth(fn components.AuthCallback) {
	s.auth = fn
}

func (s *SelfBot) IsUser(issuer string) bool {
	if s.auth == nil {
		return false
	}
	return s.auth(issuer, "")
}

func (s *SelfBot) AddFallback(fn Fallback) {
	s.messagingFallback = append([]Fallback{fn}, s.messagingFallback...)
}

func (s *SelfBot) GetCommandsList() map[string]string {
	return s.help
}

func (s *SelfBot) processNonCommandMessages(cm *chat.Message) {
	for _, fn := range s.messagingFallback {
		if fn(cm) {
			return
		}
	}
}

func (s *SelfBot) getCommandName(body string) (string, error) {
	if body[0:1] != "/" {
		return "", errors.New("not a command")
	}

	words := strings.Fields(body)

	if len(words) < 1 {
		return "", errors.New("invalid command")
	}

	return words[0][1:], nil
}
