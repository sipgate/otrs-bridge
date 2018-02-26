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
		ticket, res, body, getTicketErr := otrs.GetTicket(ticketId)
		if getTicketErr != nil {
			log.Fatal(getTicketErr)
			c.AbortWithError(http.StatusInternalServerError, getTicketErr)
		} else if res.StatusCode >= 400 || len(ticket.Ticket) == 0 {
			message := string(body[:])
			log.Println(message)
			c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": message})
		} else {
			markdownBody, listId, cardTitle := getTicketData(ticket)
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
