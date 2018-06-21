package slack

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/contract"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"net/http"
)

type TicketUpdatedInteractor struct{}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketUpdatedInteractor) HandleTicketUpdated(ticketID string, ticket otrs.Ticket, c *gin.Context) {
	c.AbortWithStatus(http.StatusTeapot)
}

func NewTicketUpdatedInteractor() contract.TicketUpdatedInteractor {
	return &TicketUpdatedInteractor{}
}
