package trello

import (
	"log"
	"strings"

	"github.com/adlio/trello"
	"github.com/lunny/html2md"
	"github.com/pkg/errors"
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

type TicketUpdatedInteractor struct{}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketUpdatedInteractor) HandleTicketUpdated(ticketID string, ticket otrs.Ticket) (contract.TicketUpdateResponse, error) {
	client := NewClient()
	card, foundCard, err := findCardByTicketID(ticketID, client)
	if err == nil {
		if foundCard {
			addCommentIfNecessary(card, ticket)
			err = processCardStateUpdate(ticket, client, card)
			if err == nil {
				return contract.SuccessfulUpdate, nil
			} else {
				return contract.UnexpectedError, err
			}
		} else {
			return contract.CardNotFound, nil
		}
	} else {
		return contract.UnexpectedError, err
	}
}

func NewTicketUpdatedInteractor() contract.TicketUpdatedInteractor {
	return &TicketUpdatedInteractor{}
}

func processCardStateUpdate(ticket otrs.Ticket, client *trello.Client, card *trello.Card) error {
	setCardPosToTicketID(card, ticket)

	if strings.Contains(ticket.State, "closed as junk") {
		return deleteCardAndRespond(client, card)
	} else if strings.Contains(ticket.State, "closed as announcement") {
		arguments := trello.Defaults()
		arguments["closed"] = "true"
		return card.Update(arguments)
	} else if strings.Contains(ticket.State, "closed") {
		listID := viper.GetString("trello.ticketDoneListID")
		return card.MoveToList(listID, trello.Defaults())
	} else if strings.Contains(ticket.State, "open") {
		listID := viper.GetString("trello.ticketDoingListID")
		return card.MoveToList(listID, trello.Defaults())
	}

	return nil
}
func setCardPosToTicketID(card *trello.Card, ticket otrs.Ticket) {
	arguments := trello.Defaults()
	arguments["pos"] = ticket.TicketID
	card.Update(arguments)
}

func addCommentIfNecessary(card *trello.Card, ticket otrs.Ticket) {
	commentCount := card.Badges.Comments
	articleCount := len(ticket.Article)
	if articleCount-commentCount > 1 {
		_, err := card.AddComment(html2md.Convert(ticket.Article[articleCount-1].Body), trello.Defaults())
		if err != nil {
			log.Println(errors.Wrap(err, "Could not add comment"))
		}
	}
}

func deleteCardAndRespond(client *trello.Client, card *trello.Card) error {
	type deleteResponse struct{}
	var response deleteResponse
	err := client.Delete("/cards/"+card.ID, trello.Defaults(), &response)
	return err
}

func findCardByTicketID(ticketID string, client *trello.Client) (*trello.Card, bool, error) {
	boardID := viper.GetString("trello.boardID")
	board, err := client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		log.Println("could not get board "+boardID, err)
		return nil, false, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Println("could not get cards", err)
		return nil, false, err
	}

	firstFoundCard := funk.Find(cards, func(card *trello.Card) bool {
		return strings.Contains(card.Name, "[#"+ticketID+"]")
	})

	if firstFoundCard != nil {
		return firstFoundCard.(*trello.Card), true, nil
	}

	return nil, false, nil
}
