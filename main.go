package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sipgate/otrs-trello-bridge/utils"
	"github.com/sipgate/otrs-trello-bridge/trello"
	"github.com/sipgate/otrs-trello-bridge/usecase"
)

func main() {
	utils.ReadConfig()
	r := gin.Default()
	trelloTicketCreatedInteractor := trello.NewTicketCreatedInteractor()
	trelloTicketCreated := usecase.NewTicketCreatedUseCase(trelloTicketCreatedInteractor)
	trelloTicketUpdatedInteractor := trello.NewTicketUpdatedInteractor()
	trelloTicketUpdated := usecase.NewTicketStateUpdatedUseCase(trelloTicketUpdatedInteractor)
	r.POST("/trello/TicketCreate/:ticketID", trelloTicketCreated.TicketCreated())
	r.POST("/trello/UpdateAllCards", trello.UpdateAllCardsHandler())
	r.POST("/trello/TicketStateUpdate/:ticketID", trelloTicketUpdated.TicketStateUpdated())



	r.Run() // listen and serve on 0.0.0.0:8080
}
