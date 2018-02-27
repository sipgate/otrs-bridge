package handler

import (
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetTicketAndHandleFailure(ticketId string, c *gin.Context) (otrs.Ticket, bool) {
	ticket, res, body, getTicketErr := otrs.GetTicket(ticketId)
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

	return ticket.Ticket[0], true
}