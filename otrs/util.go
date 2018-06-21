package otrs

import (
	"regexp"
	"log"
	"github.com/pkg/errors"
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
