package handler

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"log"
)

type TicketDeleteEvent struct {

}

func TicketDeleteHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		log.Println(buf.String())
	}
}