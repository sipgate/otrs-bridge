package contract

import (
	"github.com/sipgate/otrs-bridge/otrs"
)

type TicketUpdatedInteractor interface {
	HandleTicketUpdated(ticketID string, ticket otrs.Ticket) (TicketUpdateResponse, error)
}
