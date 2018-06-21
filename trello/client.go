package trello

import (
	"github.com/adlio/trello"
	"github.com/sipgate/otrs-bridge/utils"
	"github.com/spf13/viper"
)

// NewClient creates a new trello client with config values provided via viper config
func NewClient() *trello.Client {
	client := trello.NewClient(viper.GetString("trello.appKey"), viper.GetString("trello.token"))
	proxy := viper.GetString("trello.proxy")
	if proxy != "" {
		client.Client = utils.NewHttpClient(proxy)
	}

	return client
}
