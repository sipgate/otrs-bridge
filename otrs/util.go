package otrs

import (
	"regexp"
	"log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ExtractTicketID(name string) (string, error) {
	re, err := regexp.Compile(`^\[#(\d+)].*`)
	res := re.FindAllStringSubmatch(name, 1)

	if err != nil {
		log.Println(errors.Wrap(err, "could not compile regexp pattern"))
		return "", err
	}

	if len(res) == 1 {
		return res[0][1], nil
	}

	return "", errors.New("could not extract ticketID from '" + name + "'")
}

func MakeTicketUrl(ticket Ticket) string {
	otrsBaseURL := viper.GetString("otrs.baseUrl")
	originalTicketURL := otrsBaseURL + "/index.pl?Action=AgentTicketZoom;TicketID=" + ticket.TicketID
	return originalTicketURL
}
