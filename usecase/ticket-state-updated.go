package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
	"net/http"
)

type TicketStateUpdatedUseCase struct {
	ticketUpdatedInteractor contract.TicketUpdatedInteractor
}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketStateUpdatedUseCase) TicketStateUpdated() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketID := c.Param("ticketID")
		ticket, ok, err := otrs.GetTicketAndHandleFailure(ticketID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else if ok {
			response, err := t.ticketUpdatedInteractor.HandleTicketUpdated(ticketID, ticket)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else if response == contract.CardNotFound {
				c.AbortWithStatus(http.StatusTeapot)
			} else {
				c.AbortWithStatus(http.StatusAccepted)
			}
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

func NewTicketStateUpdatedUseCase(ticketUpdatedInteractor contract.TicketUpdatedInteractor) contract.TicketStateUpdatedUseCase {
	return &TicketStateUpdatedUseCase{ticketUpdatedInteractor}
}
