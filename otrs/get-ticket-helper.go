package otrs

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"log"
)

// GetTicketAndHandleFailure tries to get a Ticket from otrs, otherwise abort response with error
func GetTicketAndHandleFailure(ticketID string) (Ticket, bool, error) {
	ticket, res, body, getTicketErr := GetTicket(ticketID)
	if getTicketErr != nil {
		log.Println(getTicketErr)
		return Ticket{}, false, getTicketErr
	} else if res.StatusCode >= 400 || len(ticket.Ticket) == 0 {
		message := string(body[:])
		log.Println(message)
		return Ticket{}, false, errors.New(message)
	}

	queues := viper.GetStringSlice("otrs.queues")

	if len(queues) == 0 {
		return ticket.Ticket[0], true, nil
	} else if funk.ContainsString(queues, ticket.Ticket[0].Queue) {
		return ticket.Ticket[0], true, nil
	}

	return ticket.Ticket[0], false, nil
}
