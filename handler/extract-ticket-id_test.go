package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractticketID(t *testing.T) {
	ticketID, err := extractTicketID("[#12345] some ticket name")
	assert.NoError(t, err, "extract does not return an error")
	assert.Equal(t, "12345", ticketID, "extracts the correct ticketID")
}
