package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
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

func NewHttpClient(proxy string) *http.Client {
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Panicln(errors.Wrapf(err, "could not parse proxy url %s", proxy))
		}
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}

	return nil
}
