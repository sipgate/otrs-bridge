package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/adlio/trello"
	"github.com/lunny/html2md"
	"github.com/spf13/viper"
	"log"
	"fmt"
)

func TicketCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, err := otrs.GetTicket(ticketId)
		if err != nil {
			log.Fatal(err)
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			if len(ticket.Ticket) > 0 {
				markdownBody, listId, cardTitle := getTicketData(ticket)
				err := createTrelloCard(cardTitle, markdownBody, listId)

				if err == nil {
					c.AbortWithStatus(http.StatusAccepted)
				} else {
					log.Println(err)
					c.AbortWithError(500, err)
				}
			} else {
				c.AbortWithStatusJSON(404, gin.H{"message": "Ticket with ID " + ticketId + " not found"})
			}
		}
	}
}
func createTrelloCard(cardTitle string, markdownBody string, listId string) error {
	card := trello.Card{
		Name:   cardTitle,
		Desc:   markdownBody,
		IDList: listId,
	}
	client := trelloClient.NewClient()
	err := client.CreateCard(&card, trello.Defaults())
	return err
}
func getTicketData(ticket otrs.TicketResponse) (string, string, string) {
	firstTicket := ticket.Ticket[0]
	originalTicketUrl := "***Original ticket***: http://tickets.sipgate.net/otrs/index.pl?Action=AgentTicketZoom;TicketID=" + firstTicket.TicketID
	markdownBody := html2md.Convert(firstTicket.Article[0].Body)
	markdownBody = originalTicketUrl + "\n\n---\n\n" + markdownBody
	listId := viper.GetString("trello.ticketCreateListId")
	cardTitle := fmt.Sprintf("[#%s] %s", firstTicket.TicketID, firstTicket.Title)
	return markdownBody, listId, cardTitle
}
