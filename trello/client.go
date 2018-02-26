package trello

import (
	"github.com/derveloper/trello"
	"github.com/spf13/viper"
	"net/url"
	"net/http"
)

func NewClient() *trello.Client {
	client := trello.NewClient(viper.GetString("trello.appKey"), viper.GetString("trello.token"))
	proxyUrl, _ := url.Parse("http://proxy.sipgate.net:8888")
	client.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	return client
}