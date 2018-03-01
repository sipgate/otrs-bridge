package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/handler"
	"github.com/sipgate/otrs-trello-bridge/utils"
)

func main() {
	utils.ReadConfig()
	r := gin.Default()
	r.POST("/TicketCreate/:TicketId", handler.TicketCreateHandler())
	r.POST("/UpdateAllCards", handler.UpdateAllCardsHandler())
	r.POST("/TicketStateUpdate/:TicketId", handler.TicketStateUpdateHandler())
	r.Run() // listen and serve on 0.0.0.0:8080
}
