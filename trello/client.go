package trello

import (
	"github.com/adlio/trello"
	"github.com/spf13/viper"
)

func NewClient() *trello.Client {
	client := trello.NewClient(viper.GetString("trello.appKey"), viper.GetString("trello.token"))
	return client
}