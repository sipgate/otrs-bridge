package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTicketId(t *testing.T) {
	ticketId, err := extractTicketId("[#12345] some ticket name")
	assert.NoError(t, err, "extract does not return an error")
	assert.Equal(t, "12345", ticketId, "extracts the correct ticketId")
}
