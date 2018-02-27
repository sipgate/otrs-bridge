package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	trelloClient "github.sipgate.net/sipgate/otrs-trello-bride/trello"
	"strings"
	"github.com/derveloper/trello"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

type TicketStateUpdateEvent struct {
	OldTicketData struct {
		State                  string `json:"State"`
		Type                   string `json:"Type"`
		Changed                string `json:"Changed"`
		Responsible            string `json:"Responsible"`
		EscalationTime         string `json:"EscalationTime"`
		PriorityID             string `json:"PriorityID"`
		CustomerID             string `json:"CustomerID"`
		ServiceID              string `json:"ServiceID"`
		EscalationResponseTime string `json:"EscalationResponseTime"`
		Age                    int    `json:"Age"`
		CustomerUserID         string `json:"CustomerUserID"`
		UntilTime              int    `json:"UntilTime"`
		EscalationUpdateTime   string `json:"EscalationUpdateTime"`
		Lock                   string `json:"Lock"`
		ChangeBy               string `json:"ChangeBy"`
		TicketNumber           string `json:"TicketNumber"`
		StateID                string `json:"StateID"`
		Owner                  string `json:"Owner"`
		UnlockTimeout          string `json:"UnlockTimeout"`
		Title                  string `json:"Title"`
		OwnerID                string `json:"OwnerID"`
		SLAID                  string `json:"SLAID"`
		ArchiveFlag            string `json:"ArchiveFlag"`
		Priority               string `json:"Priority"`
		LockID                 string `json:"LockID"`
		TicketID               string `json:"TicketID"`
		TypeID                 string `json:"TypeID"`
		RealTillTimeNotUsed    string `json:"RealTillTimeNotUsed"`
		StateType              string `json:"StateType"`
		EscalationSolutionTime string `json:"EscalationSolutionTime"`
		CreateTimeUnix         string `json:"CreateTimeUnix"`
		ResponsibleID          string `json:"ResponsibleID"`
		QueueID                string `json:"QueueID"`
		Created                string `json:"Created"`
		GroupID                string `json:"GroupID"`
		Queue                  string `json:"Queue"`
		CreateBy               string `json:"CreateBy"`
	} `json:"OldTicketData"`
}

func TicketStateUpdateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		ticketId := c.Param("TicketId")
		ticket, ok := GetTicketAndHandleFailure(ticketId, c)
		if ok {
			if strings.Contains(ticket.State, "closed") {
				client := trelloClient.NewClient()
				card, err := findCard(ticketId, client)
				if err == nil {
					err := card.MoveToList(viper.GetString("trello.ticketDoneListId"), trello.Defaults())
					if err == nil {
						c.AbortWithStatus(http.StatusAccepted)
					} else {
						log.Println(err)
						c.AbortWithError(http.StatusInternalServerError, err)
					}
				} else {
					log.Println(err)
					c.AbortWithError(http.StatusInternalServerError, err)
				}
			} else {
				c.AbortWithStatus(http.StatusTeapot)
			}
		}
	}
}
func findCard(ticketId string, client *trello.Client) (*trello.Card, error) {
	boardId := viper.GetString("trello.boardId")
	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Println("could not get board " + boardId, err)
		return nil, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Println("could not get cards", err)
		return nil, err
	}

	firstFoundCard := funk.Find(cards, func(card *trello.Card) bool {
		return strings.Contains(card.Name, "[#"+ticketId+"]")
	})

	return firstFoundCard.(*trello.Card), nil
}