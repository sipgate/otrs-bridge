package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/derveloper/trello"
	"github.com/gin-gonic/gin"
	"github.com/lunny/html2md"
	"github.com/pkg/errors"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	trelloClient "github.com/sipgate/otrs-trello-bridge/trello"
	"github.com/sipgate/otrs-trello-bridge/utils"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// TicketStateUpdateHandler handles ticket state update events from otrs
func TicketStateUpdateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketID := c.Param("ticketID")
		ticket, ok := GetTicketAndHandleFailure(ticketID, c)
		if ok {
			client := trelloClient.NewClient()
			card, foundCard, err := findCardByTicketID(ticketID, client)
			if foundCard {
				utils.DoIfNoErrorOrAbort(c, err, func() {
					addCommentIfNecessary(card, ticket)
					processCardStateUpdate(ticket, c, err, client, card)
				})
			} else {
				c.AbortWithStatus(http.StatusTeapot)
			}
		}
	}
}
func processCardStateUpdate(ticket otrs.Ticket, c *gin.Context, err error, client *trello.Client, card *trello.Card) {
	if strings.Contains(ticket.State, "closed as junk") {
		utils.DoIfNoErrorOrAbort(c, err, func() {
			deleteCardAndRespond(client, card, c)
		})
	} else if strings.Contains(ticket.State, "closed as announcement") {
		arguments := trello.Defaults()
		arguments["closed"] = "true"
		err := card.Update(arguments)
		utils.DoIfNoErrorOrAbort(c, err, func() {
			c.AbortWithStatus(http.StatusAccepted)
		})
	} else if strings.Contains(ticket.State, "closed") {
		listID := viper.GetString("trello.ticketDoneListID")
		moveCardAndRespond(card, listID, c)
	} else if strings.Contains(ticket.State, "open") {
		listID := viper.GetString("trello.ticketDoingListID")
		moveCardAndRespond(card, listID, c)
	} else {
		c.AbortWithStatus(http.StatusTeapot)
	}
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

func deleteCardAndRespond(client *trello.Client, card *trello.Card, c *gin.Context) {
	type deleteResponse struct{}
	var response deleteResponse
	err := client.Delete("/cards/"+card.ID, trello.Defaults(), &response)
	utils.DoIfNoErrorOrAbort(c, err, func() {
		c.AbortWithStatus(http.StatusAccepted)
	})
}

func moveCardAndRespond(card *trello.Card, listID string, c *gin.Context) {
	err := card.MoveToList(listID, trello.Defaults())
	utils.DoIfNoErrorOrAbort(c, err, func() {
		c.AbortWithStatus(http.StatusAccepted)
	})
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
