package handler

import (
	"github.com/gin-gonic/gin"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/spf13/viper"
	"github.com/derveloper/trello"
	"github.sipgate.net/sipgate/otrs-trello-bride/utils"
	"github.com/thoas/go-funk"
	"regexp"
	"errors"
)

func UpdateAllCardsHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		client := trelloClient.NewClient()
		board, err := client.GetBoard(viper.GetString("trello.boardId"), trello.Defaults())
		utils.DoIfNoErrorOrAbort(c, err, func() {
			cards, err := board.GetCards(trello.Defaults())
			utils.DoIfNoErrorOrAbort(c, err, func() {
				funk.Map(cards, func(card *trello.Card) string {
					return card.Name
				})
			})
		})
	}
}

func extractTicketId(name string) (string, error) {
	re, err := regexp.Compile(`^\[#(\d+)].*`)
	res := re.FindAllStringSubmatch(name, 1)

	if err != nil {
		return "", err
	}

	if len(res) == 1 {
		return res[0][1], nil
	} else {
		return "", errors.New("could not extract ticketId from '" + name + "'")
	}
}