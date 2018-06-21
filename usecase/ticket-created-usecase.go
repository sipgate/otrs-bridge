package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/contract"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"log"
	"net/http"
)

type TicketCreatedUseCase struct {
	ticketCreatedInteractor contract.TicketCreatedInteractor
}

// TicketCreateHandler handles ticket creation events from otrs
func (t *TicketCreatedUseCase) TicketCreated() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketID := c.Param("ticketID")
		ticket, ok := otrs.GetTicketAndHandleFailure(ticketID, c)
		if ok {
			err := t.ticketCreatedInteractor.HandleTicketCreated(ticketID, ticket)
			if err == nil {
				c.AbortWithStatus(http.StatusAccepted)
			} else {
				log.Println(err)
				c.AbortWithError(500, err)
			}
		}
	}
}

func NewTicketCreatedUseCase(interactor contract.TicketCreatedInteractor) contract.TicketCreatedUseCase {
	return &TicketCreatedUseCase{interactor}
}
