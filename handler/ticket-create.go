package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/derveloper/trello"
	"github.com/gin-gonic/gin"
	"github.com/lunny/html2md"
	"github.com/pkg/errors"
	"github.com/sipgate/otrs-trello-bridge/otrs"
	trelloClient "github.com/sipgate/otrs-trello-bridge/trello"
	"github.com/spf13/viper"
)

//TicketCreateHandler handles ticket creation events from otrs
func TicketCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketID := c.Param("ticketID")
		ticket, ok := GetTicketAndHandleFailure(ticketID, c)
		if ok {
			markdownBody, listID, cardTitle := getTicketDataForCard(ticket)
			client := trelloClient.NewClient()
			err := createTrelloCard(cardTitle, markdownBody, listID, client)
			card, cardFound, err := findCardByTicketID(ticketID, client)
			if cardFound {
				arguments := trello.Defaults()
				arguments["idLabels"] = viper.GetString("trello.soonLabelId")
				err := card.Update(arguments)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not label card"))
				}
			}
			if err == nil {
				c.AbortWithStatus(http.StatusAccepted)
			} else {
				log.Println(err)
				c.AbortWithError(500, err)
			}
		}
	}
}

func createTrelloCard(cardTitle string, markdownBody string, listID string, client *trello.Client) error {
	card := trello.Card{
		Name:   cardTitle,
		Desc:   markdownBody,
		IDList: listID,
	}
	err := client.CreateCard(&card, trello.Defaults())
	return err
}

func getTicketDataForCard(ticket otrs.Ticket) (string, string, string) {
	otrsBaseURL := viper.GetString("otrs.baseUrl")
	originalTicketURL := "***Original ticket***: " + otrsBaseURL + "/index.pl?Action=AgentTicketZoom;ticketID=" + ticket.ticketID
	markdownBody := html2md.Convert(ticket.Article[0].Body)
	markdownBody = originalTicketURL + "\n\n---\n\n" + markdownBody
	listID := viper.GetString("trello.ticketCreatelistID")
	cardTitle := fmt.Sprintf("[#%s] %s", ticket.ticketID, ticket.Title)
	return markdownBody, listID, cardTitle
}
