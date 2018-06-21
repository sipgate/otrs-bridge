package slack

import (
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sipgate/otrs-trello-bridge/contract"
)

type TicketUpdatedInteractor struct {}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketUpdatedInteractor) HandleTicketUpdated(ticketID string, ticket otrs.Ticket, c *gin.Context) {
	c.AbortWithStatus(http.StatusTeapot)
}

func NewTicketUpdatedInteractor() contract.TicketUpdatedInteractor {
	return &TicketUpdatedInteractor{}
}
