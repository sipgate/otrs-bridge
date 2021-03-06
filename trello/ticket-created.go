package trello

import (
	"fmt"
	"github.com/adlio/trello"
	"github.com/lunny/html2md"
	"github.com/pkg/errors"
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type TicketCreatedInteractor struct{}

// TicketCreateHandler handles ticket creation events from otrs
func (t *TicketCreatedInteractor) HandleTicketCreated(ticketID string, ticket otrs.Ticket) error {
	markdownBody, listID, cardTitle := getTicketDataForCard(ticket)
	client := NewClient()
	err := createTrelloCard(ticketID, cardTitle, markdownBody, listID, client)
	card, cardFound, err := findCardByTicketID(ticketID, client)
	if cardFound {
		arguments := trello.Defaults()
		arguments["idLabels"] = viper.GetString("trello.soonLabelId")
		err := card.Update(arguments)
		if err != nil {
			log.Println(errors.Wrap(err, "Could not label card"))
		}
	}

	return err
}

func NewTicketCreatedInteractor() contract.TicketCreatedInteractor {
	return &TicketCreatedInteractor{}
}

func createTrelloCard(ticketID string, cardTitle string, markdownBody string, listID string, client *trello.Client) error {
	cardPos, err := strconv.Atoi(ticketID)
	if err != nil {
		return err
	}
	card := trello.Card{
		Name:   cardTitle,
		Desc:   markdownBody,
		IDList: listID,
		Pos:    float64(cardPos),
	}
	return client.CreateCard(&card, trello.Defaults())
}

func getTicketDataForCard(ticket otrs.Ticket) (string, string, string) {
	originalTicketURL := otrs.MakeTicketUrl(ticket)
	markdownBody := html2md.Convert(ticket.Article[0].Body)
	markdownBody = "***Original ticket***: " + originalTicketURL + "\n\n---\n\n" + markdownBody
	listID := viper.GetString("trello.ticketCreateListID")
	cardTitle := fmt.Sprintf("[#%s] %s", ticket.TicketID, ticket.Title)
	return markdownBody, listID, cardTitle
}
