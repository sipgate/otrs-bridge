package slack

import (
	"github.com/sipgate/otrs-trello-bridge/contract"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"github.com/nlopes/slack"
	"fmt"
	"github.com/spf13/viper"
)

type TicketCreatedInteractor struct {}

func (t *TicketCreatedInteractor) HandleTicketCreated(ticketID string, ticket otrs.Ticket) error {
	api := slack.New(viper.GetString("slack.apiToken"))
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: ticket.Title,
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(viper.GetString("slack.channelId"), "New OTRS Ticket", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	return err
}

func NewTicketCreatedInteractor() contract.TicketCreatedInteractor {
	return &TicketCreatedInteractor{}
}
