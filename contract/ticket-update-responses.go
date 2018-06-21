package contract

type TicketUpdateResponse int

const (
	UnexpectedError TicketUpdateResponse = iota
	CardNotFound
	SuccessfulUpdate
)
