package contract

import "github.com/gin-gonic/gin"

type TicketCreatedUseCase interface {
	TicketCreated() func(c *gin.Context)
}
