package components

import (
	"github.com/joinself/self-go-sdk/chat"
)

type HelpComponent struct {
	header  string
	footer  string
	summary string
}

type HelpComponentConfig struct {
	Header  string
	Footer  string
	Summary string
}

func NewHelpComponent(cfg HelpComponentConfig) *HelpComponent {
	summary := "Prints contextual help for this bot"
	if len(cfg.Summary) == 0 {
		summary = cfg.Summary
	}

	header := "See below the list of available commands"
	if len(cfg.Header) == 0 {
		header = cfg.Header
	}

	return &HelpComponent{
		header:  header,
		footer:  cfg.Footer,
		summary: summary,
	}
}

func (i *HelpComponent) RecordCommands(r CommandRecorder) {
	i.registerHelpCommand(r)
}

func (i *HelpComponent) AfterStartHook(r CommandRecorder) {}

func (i *HelpComponent) registerHelpCommand(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "help",
		Summary: i.summary,
		Callback: func(cmd string, cm *chat.Message) string {
			commands := r.GetCommandsList()

			helpMsg := i.header + "\n"
			for k, v := range commands {
				helpMsg += "\n - <b>" + k + "</b> : " + v
			}
			helpMsg += "\n" + i.footer

			return helpMsg
		},
	})
}
