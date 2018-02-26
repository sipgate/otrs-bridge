package main

import (
	"github.com/gin-gonic/gin"
	"github.sipgate.net/sipgate/otrs-trello-bride/handler"
	"github.sipgate.net/sipgate/otrs-trello-bride/utils"
	"os"
)

func main() {
	os.Setenv("HTTP_PROXY", "http://proxy.netzquadrat.net:8888")
	utils.ReadConfig()
	r := gin.Default()
	r.POST("/TicketCreate/:TicketId", handler.TicketCreateHandler())
	r.POST("/TicketDelete/:TicketId", handler.TicketDeleteHandler())
	r.POST("/TicketStateUpdate/:TicketId", handler.TicketStateUpdateHandler())
	r.Run() // listen and serve on 0.0.0.0:8080
}
