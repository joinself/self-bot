package components

import (
	"encoding/json"
	"fmt"
	"log"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/chat"
)

type InvitesSummaryComponent struct {
	client *selfsdk.Client
}

type InvitesSummaryComponentConfig struct {
	Client *selfsdk.Client
}

func NewInvitesSummaryComponent(cfg InvitesSummaryComponentConfig) *InvitesSummaryComponent {
	return &InvitesSummaryComponent{
		client: cfg.Client,
	}
}

func (i *InvitesSummaryComponent) RecordCommands(r CommandRecorder) {
	i.registerSummaryCommand(r)
}

func (i *InvitesSummaryComponent) AfterStartHook(r CommandRecorder) {}

func (i *InvitesSummaryComponent) registerSummaryCommand(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "summary",
		Summary: "Displays a summary for your signup incentive program.",
		Callback: func(cmd string, cm *chat.Message) string {
			resp, err := i.client.Rest().Get("/v1/invites")
			if err != nil {
				log.Println(err.Error())
				return "An internal error has ocurred"
			}

			type summary struct {
				Remaining int `json:"remaining"`
				Points    int `json:"points"`
			}
			var s summary
			json.Unmarshal(resp, &s)

			return fmt.Sprintf("You’ve earned %d points by inviting people to join Self! \n\nKeep spreading the word and earning points", s.Points)
		},
	})
}
