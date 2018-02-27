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
			card, foundCard, err := findCard(ticketId, client)
			if foundCard {
				utils.DoIfNoErrorOrAbort(c, err, func() {
					if strings.Contains(ticket.State, "closed as junk") {
						utils.DoIfNoErrorOrAbort(c, err, func() {
							DeleteCardAndRespond(client, card, c)
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
			} else {
				c.AbortWithStatus(http.StatusTeapot)
			}
		}
	}
}
func DeleteCardAndRespond(client *trello.Client, card *trello.Card, c *gin.Context) {
	type deleteResponse struct{}
	var response deleteResponse
	err := client.Delete("/cards/"+card.ID, trello.Defaults(), &response)
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

func findCard(ticketId string, client *trello.Client) (*trello.Card, bool, error) {
	boardId := viper.GetString("trello.boardId")
	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Println("could not get board "+boardId, err)
		return nil, false, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Println("could not get cards", err)
		return nil, false, err
	}

	firstFoundCard := funk.Find(cards, func(card *trello.Card) bool {
		return strings.Contains(card.Name, "[#"+ticketId+"]")
	})

	if firstFoundCard != nil {
		return firstFoundCard.(*trello.Card), true, nil
	}

	return nil, false, nil
}
