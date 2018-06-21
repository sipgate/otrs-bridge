package slack

import (
	"github.com/sipgate/otrs-bridge/contract"
	"github.com/sipgate/otrs-bridge/otrs"
)

type TicketUpdatedInteractor struct{}

// TicketStateUpdateHandler handles ticket state update events from otrs
func (t *TicketUpdatedInteractor) HandleTicketUpdated(ticketID string, ticket otrs.Ticket) (contract.TicketUpdateResponse, error) {
	return contract.SuccessfulUpdate, nil
}

func NewTicketUpdatedInteractor() contract.TicketUpdatedInteractor {
	return &TicketUpdatedInteractor{}
}
