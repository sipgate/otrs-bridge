package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
)

type TicketStateUpdatedUseCase struct {
	ticketUpdatedInteractor contract.TicketUpdatedInteractor
}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketStateUpdatedUseCase) TicketStateUpdated() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketID := c.Param("ticketID")
		ticket, ok := otrs.GetTicketAndHandleFailure(ticketID, c)
		if ok {
			t.ticketUpdatedInteractor.HandleTicketUpdated(ticketID, ticket, c)
		}
	}
}

func NewTicketStateUpdatedUseCase(ticketUpdatedInteractor contract.TicketUpdatedInteractor) contract.TicketStateUpdatedUseCase {
	return &TicketStateUpdatedUseCase{ticketUpdatedInteractor}
}
