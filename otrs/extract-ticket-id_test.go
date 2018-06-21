package otrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTicketID(t *testing.T) {
	ticketID, err := ExtractTicketID("[#12345] some ticket name")
	assert.NoError(t, err, "extract does not return an error")
	assert.Equal(t, "12345", ticketID, "extracts the correct ticketID")
}
