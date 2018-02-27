package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DoIfNoErrorOrAbort(c *gin.Context, err error, f func())  {
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		f()
	}
}