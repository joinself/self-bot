package components

import (
	"strings"

	"github.com/joinself/self-go-sdk/chat"
)

type AuthCallback func(issuer, cmd string) bool

type CommandRecorder interface {
	RecordCommand(Command)
	SetAuth(AuthCallback)
	GetCommandsList() map[string]string
	IsUser(issuer string) bool
}

type CommandComponent interface {
	RecordCommands(r CommandRecorder)
	AfterStartHook(r CommandRecorder)
}

type CommandCallback func(string, *chat.Message) string

type Command struct {
	Name     string
	Summary  string
	Callback CommandCallback
}

func GetCommandParams(body string) []string {
	words := strings.Fields(body)
	if len(words) <= 1 {
		return []string{}
	}

	return words[1:]
}
