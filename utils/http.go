package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DoIfNoErrorOrAbort(c *gin.Context, err error, f func()) {
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		f()
	}
}
