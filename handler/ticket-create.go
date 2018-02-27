package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/derveloper/trello"
	"github.com/lunny/html2md"
	"github.com/spf13/viper"
	"log"
	"fmt"
)

func TicketCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, ok := GetTicketAndHandleFailure(ticketId, c)
		if ok {
			markdownBody, listId, cardTitle := getTicketDataForCard(ticket)
			err := createTrelloCard(cardTitle, markdownBody, listId)

			if err == nil {
				c.AbortWithStatus(http.StatusAccepted)
			} else {
				log.Println(err)
				c.AbortWithError(500, err)
			}
		}
	}
}

func createTrelloCard(cardTitle string, markdownBody string, listId string) error {
	card := trello.Card{
		Name:   cardTitle,
		Desc:   markdownBody,
		IDList: listId,
		Labels: []*trello.Label {
			{Name: "bald"},
		},
	}
	client := trelloClient.NewClient()
	err := client.CreateCard(&card, trello.Defaults())
	return err
}

func getTicketDataForCard(ticket otrs.Ticket) (string, string, string) {
	originalTicketUrl := "***Original ticket***: http://tickets.sipgate.net/otrs/index.pl?Action=AgentTicketZoom;TicketID=" + ticket.TicketID
	markdownBody := html2md.Convert(ticket.Article[0].Body)
	markdownBody = originalTicketUrl + "\n\n---\n\n" + markdownBody
	listId := viper.GetString("trello.ticketCreateListId")
	cardTitle := fmt.Sprintf("[#%s] %s", ticket.TicketID, ticket.Title)
	return markdownBody, listId, cardTitle
}
