package spacer

import (
	"regexp"
)

var regex = regexp.MustCompile(`.*\[.+\].*`)

func Run(s string) bool {
	return regex.MatchString(s)
}
