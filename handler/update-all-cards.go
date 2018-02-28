package handler

import (
	"github.com/gin-gonic/gin"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/spf13/viper"
	"github.com/derveloper/trello"
	"github.sipgate.net/sipgate/otrs-trello-bride/utils"
	"github.com/thoas/go-funk"
	"regexp"
	"log"
	"github.com/pkg/errors"
	"net/http"
)

func UpdateAllCardsHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		client := trelloClient.NewClient()
		board, err := client.GetBoard(viper.GetString("trello.boardId"), trello.Defaults())
		utils.DoIfNoErrorOrAbort(c, err, func() {
			cards, err := board.GetCards(trello.Defaults())
			utils.DoIfNoErrorOrAbort(c, err, func() {
				ticketIds := funk.Map(cards, func(card *trello.Card) string {
					id, _ := extractTicketId(card.Name)
					return id
				}).([]string)

				ticketIds = funk.FilterString(ticketIds, func(id string) bool {
					return id != ""
				})

				for _, id := range ticketIds {
					_, err := http.Post("http://localhost:8080/TicketStateUpdate/" + id, "application/json", nil)
					if err == nil {
						log.Println("updated card for ticket " + id)
					} else {
						log.Println(errors.Wrapf(err, "could not update card for ticket %s", id))
					}
				}

				c.AbortWithStatus(http.StatusAccepted)
			})
		})
	}
}

func extractTicketId(name string) (string, error) {
	re, err := regexp.Compile(`^\[#(\d+)].*`)
	res := re.FindAllStringSubmatch(name, 1)

	if err != nil {
		log.Println(errors.Wrap(err, "could not compile regexp pattern"))
		return "", err
	}

	if len(res) == 1 {
		return res[0][1], nil
	} else {
		return "", errors.New("could not extract ticketId from '" + name + "'")
	}
}