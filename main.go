package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/slack"
	"github.com/sipgate/otrs-trello-bridge/trello"
	"github.com/sipgate/otrs-trello-bridge/usecase"
	"github.com/sipgate/otrs-trello-bridge/utils"
	"github.com/spf13/viper"
	"log"
)

func main() {
	utils.ReadConfig()

	hasFunctionality := viper.GetBool("trello.enabled") && viper.GetBool("slack.enabled")

	if !hasFunctionality {
		log.Fatal("no integrations are enabled, exiting")
	}

	r := gin.Default()

	if viper.GetBool("trello.enabled") {
		initTrelloRoutes(r)
	} else {
		log.Println("trello support is disabled")
	}

	if viper.GetBool("slack.enabled") {
		initSlackRoutes(r)
	} else {
		log.Println("slack support is disabled")
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

func initTrelloRoutes(r *gin.Engine) {
	ticketCreatedInteractor := trello.NewTicketCreatedInteractor()
	ticketCreated := usecase.NewTicketCreatedUseCase(ticketCreatedInteractor)
	r.POST("/trello/TicketCreate/:ticketID", ticketCreated.TicketCreated())

	ticketUpdatedInteractor := trello.NewTicketUpdatedInteractor()
	ticketUpdated := usecase.NewTicketStateUpdatedUseCase(ticketUpdatedInteractor)
	r.POST("/trello/TicketStateUpdate/:ticketID", ticketUpdated.TicketStateUpdated())

	r.POST("/trello/UpdateAllCards", trello.UpdateAllCardsHandler())
}

func initSlackRoutes(r *gin.Engine) {
	ticketCreatedInteractor := slack.NewTicketCreatedInteractor()
	ticketCreated := usecase.NewTicketCreatedUseCase(ticketCreatedInteractor)
	r.POST("/slack/TicketCreate/:ticketID", ticketCreated.TicketCreated())

	ticketUpdatedInteractor := slack.NewTicketUpdatedInteractor()
	ticketUpdated := usecase.NewTicketStateUpdatedUseCase(ticketUpdatedInteractor)
	r.POST("/slack/TicketStateUpdate/:ticketID", ticketUpdated.TicketStateUpdated())
}
