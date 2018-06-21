package contract

import (
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"github.com/gin-gonic/gin"
)

type TicketUpdatedInteractor interface {
	HandleTicketUpdated(ticketID string, ticket otrs.Ticket, c *gin.Context)
}
