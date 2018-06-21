package trello

import (
	"github.com/derveloper/trello"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	"github.com/sipgate/otrs-trello-bridge/utils"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"log"
	"net/http"
)

// UpdateAllCardsHandler calls TicketStateUpdate for all cards in board
func UpdateAllCardsHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		client := NewClient()
		board, err := client.GetBoard(viper.GetString("trello.boardID"), trello.Defaults())
		utils.DoIfNoErrorOrAbort(c, err, func() {
			cards, err := board.GetCards(trello.Defaults())
			utils.DoIfNoErrorOrAbort(c, err, func() {
				ticketIDs := funk.Map(cards, func(card *trello.Card) string {
					id, _ := otrs.ExtractTicketID(card.Name)
					return id
				}).([]string)

				ticketIDs = funk.FilterString(ticketIDs, func(id string) bool {
					return id != ""
				})

				for _, id := range ticketIDs {
					_, err := http.Post("http://localhost:8080/TicketStateUpdate/"+id, "application/json", nil)
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
