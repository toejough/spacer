package spacer

import "strings"

func Run(s string) bool {
	return strings.Contains(s, "[")
}
