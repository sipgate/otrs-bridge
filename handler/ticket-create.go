package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.sipgate.net/sipgate/otrs-trello-bride/otrs"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"github.com/derveloper/trello"
	"github.com/lunny/html2md"
	"github.com/spf13/viper"
	"log"
	"fmt"
	"github.com/pkg/errors"
)

func TicketCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, ok := GetTicketAndHandleFailure(ticketId, c)
		if ok {
			markdownBody, listId, cardTitle := getTicketDataForCard(ticket)
			client := trelloClient.NewClient()
			err := createTrelloCard(cardTitle, markdownBody, listId, client)
			card, cardFound, err := findCardByTicketId(ticketId, client)
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

func createTrelloCard(cardTitle string, markdownBody string, listId string, client *trello.Client) error {
	card := trello.Card{
		Name:   cardTitle,
		Desc:   markdownBody,
		IDList: listId,
	}
	err := client.CreateCard(&card, trello.Defaults())
	return err
}

func getTicketDataForCard(ticket otrs.Ticket) (string, string, string) {
	otrsBaseUrl := viper.GetString("otrs.baseUrl")
	originalTicketUrl := "***Original ticket***: " + otrsBaseUrl + "/index.pl?Action=AgentTicketZoom;TicketID=" + ticket.TicketID
	markdownBody := html2md.Convert(ticket.Article[0].Body)
	markdownBody = originalTicketUrl + "\n\n---\n\n" + markdownBody
	listId := viper.GetString("trello.ticketCreateListId")
	cardTitle := fmt.Sprintf("[#%s] %s", ticket.TicketID, ticket.Title)
	return markdownBody, listId, cardTitle
}
