package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// GetTicketAndHandleFailure tries to get a Ticket from otrs, otherwise abort response with error
func GetTicketAndHandleFailure(ticketID string, c *gin.Context) (otrs.Ticket, bool) {
	ticket, res, body, getTicketErr := otrs.GetTicket(ticketID)
	if getTicketErr != nil {
		log.Fatal(getTicketErr)
		c.AbortWithError(http.StatusInternalServerError, getTicketErr)
		return otrs.Ticket{}, false
	} else if res.StatusCode >= 400 || len(ticket.Ticket) == 0 {
		message := string(body[:])
		log.Println(message)
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": message})
		return otrs.Ticket{}, false
	}

	queues := viper.GetStringSlice("otrs.queues")

	if len(queues) == 0 {
		return ticket.Ticket[0], true
	} else if funk.ContainsString(queues, ticket.Ticket[0].Queue) {
		return ticket.Ticket[0], true
	}

	return ticket.Ticket[0], false
}
