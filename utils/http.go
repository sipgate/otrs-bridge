package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DoIfNoErrorOrAbort wraps an error check and executes the given function otherwise it aborts the request with InternalServerError
func DoIfNoErrorOrAbort(c *gin.Context, err error, f func()) {
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		f()
	}
}
