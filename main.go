package main

import (
	"github.com/gin-gonic/gin"
	"github.sipgate.net/sipgate/otrs-trello-bride/handler"
	"github.sipgate.net/sipgate/otrs-trello-bride/utils"
)

func main() {
	utils.ReadConfig()
	r := gin.Default()
	r.POST("/TicketCreate/:TicketId", handler.TicketCreateHandler())
	r.POST("/TicketDelete/:TicketId", handler.TicketDeleteHandler())
	r.POST("/TicketStateUpdate/:TicketId", handler.TicketCreateHandler())
	r.Run() // listen and serve on 0.0.0.0:8080
}
