package contract

import "github.com/gin-gonic/gin"

type TicketStateUpdatedUseCase interface {
	TicketStateUpdated() func(c *gin.Context)
}
