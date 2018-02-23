package utils

import (
	"github.com/lunny/html2md"
)

func convert(html string) string {
	return html2md.Convert(html)
}