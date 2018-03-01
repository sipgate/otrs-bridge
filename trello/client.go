package trello

import (
	"log"
	"net/http"
	"net/url"

	"github.com/derveloper/trello"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewClient() *trello.Client {
	client := trello.NewClient(viper.GetString("trello.appKey"), viper.GetString("trello.token"))
	proxy := viper.GetString("trello.proxy")
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Panicln(errors.Wrapf(err, "could not parse proxy url %s", proxy))
		}
		client.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	return client
}
