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
			firstTicket := ticket.Ticket[0]
			markdownBody := html2md.Convert(firstTicket.Article[0].Body)

			listId := viper.GetString("trello.ticketCreateListId")
			cardTitle := fmt.Sprintf("[#%s] %s", firstTicket.TicketID, firstTicket.Title)
			card := trello.Card{Name: cardTitle, Desc: markdownBody, IDList: listId}

			client := trelloClient.NewClient()
			err := client.CreateCard(&card, trello.Defaults())

			if err == nil {
				c.AbortWithStatus(http.StatusAccepted)
			} else {
				log.Println(err)
				c.AbortWithError(500, err)
			}
		}
	}
}
