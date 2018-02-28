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
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	"github.com/lunny/html2md"
	"github.com/pkg/errors"
)

func TicketStateUpdateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, ok := GetTicketAndHandleFailure(ticketId, c)
		if ok {
			client := trelloClient.NewClient()
			card, foundCard, err := findCardByTicketId(ticketId, client)
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
			DeleteCardAndRespond(client, card, c)
		})
	} else if strings.Contains(ticket.State, "closed as announcement") {
		arguments := trello.Defaults()
		arguments["closed"] = "true"
		err := card.Update(arguments)
		utils.DoIfNoErrorOrAbort(c, err, func() {
			c.AbortWithStatus(http.StatusAccepted)
		})
	} else if strings.Contains(ticket.State, "closed") {
		listId := viper.GetString("trello.ticketDoneListId")
		moveCardAndRespond(card, listId, c)
	} else if strings.Contains(ticket.State, "open") {
		listId := viper.GetString("trello.ticketDoingListId")
		moveCardAndRespond(card, listId, c)
	} else {
		c.AbortWithStatus(http.StatusTeapot)
	}
}

func addCommentIfNecessary(card *trello.Card, ticket otrs.Ticket) {
	commentCount := card.Badges.Comments
	articleCount := len(ticket.Article)
	if articleCount - commentCount > 1 {
		_, err := card.AddComment(html2md.Convert(ticket.Article[articleCount - 1].Body), trello.Defaults())
		if err != nil {
			log.Println(errors.Wrap(err, "Could not add comment"))
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

func findCardByTicketId(ticketId string, client *trello.Client) (*trello.Card, bool, error) {
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
