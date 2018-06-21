package otrs

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// GetTicketAndHandleFailure tries to get a Ticket from otrs, otherwise abort response with error
func GetTicketAndHandleFailure(ticketID string, c *gin.Context) (Ticket, bool) {
	ticket, res, body, getTicketErr := GetTicket(ticketID)
	if getTicketErr != nil {
		log.Fatal(getTicketErr)
		c.AbortWithError(http.StatusInternalServerError, getTicketErr)
		return Ticket{}, false
	} else if res.StatusCode >= 400 || len(ticket.Ticket) == 0 {
		message := string(body[:])
		log.Println(message)
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": message})
		return Ticket{}, false
	}

	queues := viper.GetStringSlice("otrs.queues")

	if len(queues) == 0 {
		return ticket.Ticket[0], true
	} else if funk.ContainsString(queues, ticket.Ticket[0].Queue) {
		return ticket.Ticket[0], true
	}

	return ticket.Ticket[0], false
}
