package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/otrs"
)

type TicketUpdatedInteractor interface {
	HandleTicketUpdated(ticketID string, ticket otrs.Ticket, c *gin.Context)
}
