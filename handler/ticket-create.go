package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	"log"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/adlio/trello"
	"github.com/lunny/html2md"
	"github.com/spf13/viper"
)

func TicketCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, err := otrs.GetTicket(ticketId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			log.Println(ticket)
			client := trelloClient.NewClient()
			list, err := client.GetList(viper.GetString("trello.ticketCreateListId"), trello.Defaults())
			if err != nil {
				firstTicket := ticket.Ticket[0]
				list.AddCard(
					&trello.Card{Name: firstTicket.Title, Desc: html2md.Convert(firstTicket.Article[0].Body)},
					trello.Defaults())
				c.AbortWithStatus(http.StatusAccepted)
			} else {
				c.AbortWithError(500, err)
			}
		}
	}
}
