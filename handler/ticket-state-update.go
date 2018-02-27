package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"strings"
	"github.com/derveloper/trello"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"github.sipgate.net/sipgate/otrs-trello-bride/utils"
)

func TicketStateUpdateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, ok := GetTicketAndHandleFailure(ticketId, c)
		if ok {
			client := trelloClient.NewClient()
			card, err := findCard(ticketId, client)
			utils.DoIfNoErrorOrAbort(c, err, func() {
				if strings.Contains(ticket.State, "closed as junk") {
					utils.DoIfNoErrorOrAbort(c, err, func() {
						DeleteCard(client, card, c)
					})
				} else if strings.Contains(ticket.State, "closed as announcement") {
					card.Closed = true
					err := card.Update(trello.Defaults())
					utils.DoIfNoErrorOrAbort(c, err, func() {
						c.AbortWithStatus(http.StatusAccepted)
					})
				} else if strings.Contains(ticket.State, "closed") {
					listId := viper.GetString("trello.ticketDoneListId")
					moveCardAndRespond(card, listId, c)
				} else {
					c.AbortWithStatus(http.StatusTeapot)
				}
			})
		}
	}
}
func DeleteCard(client *trello.Client, card *trello.Card, c *gin.Context) {
	err := client.Delete("/1/cards/"+card.ID, trello.Defaults(), trello.Card{})
	utils.DoIfNoErrorOrAbort(c, err, func() {
		c.AbortWithStatus(http.StatusAccepted)
	})
}

func moveCardAndRespond(card *trello.Card, listId string, c *gin.Context) {
	err := card.MoveToList(listId, trello.Defaults())
	utils.DoIfNoErrorOrAbort(c, err, func() {
		c.AbortWithStatus(http.StatusAccepted)
	})
}

func findCard(ticketId string, client *trello.Client) (*trello.Card, error) {
	boardId := viper.GetString("trello.boardId")
	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Println("could not get board " + boardId, err)
		return nil, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Println("could not get cards", err)
		return nil, err
	}

	firstFoundCard := funk.Find(cards, func(card *trello.Card) bool {
		return strings.Contains(card.Name, "[#"+ticketId+"]")
	})

	return firstFoundCard.(*trello.Card), nil
}
