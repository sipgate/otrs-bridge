package contract

import "github.com/sipgate/otrs-bridge/otrs"

type TicketCreatedInteractor interface {
	HandleTicketCreated(ticketID string, ticket otrs.Ticket) error
}
