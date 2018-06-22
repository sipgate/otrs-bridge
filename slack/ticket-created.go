package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
	"github.com/sipgate/otrs-bridge/utils"
	"github.com/spf13/viper"
)

type TicketCreatedInteractor struct{}

func (t *TicketCreatedInteractor) HandleTicketCreated(ticketID string, ticket otrs.Ticket) error {
	proxy := viper.GetString("slack.proxy")
	var api *slack.Client

	if proxy != "" {
		client := slack.OptionHTTPClient(utils.NewHttpClient(proxy))
		api = slack.New(viper.GetString("slack.apiToken"), client)
	} else {
		api = slack.New(viper.GetString("slack.apiToken"))
	}

	params := slack.PostMessageParameters{}
	params.Username = "otrs-bridge"
	attachment := slack.Attachment{
		Title:     ticket.Title,
		TitleLink: otrs.MakeTicketUrl(ticket),
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
