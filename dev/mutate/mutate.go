// Package mutate provides mutation testing functionality.
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	// Load the kinds of mutations to test
	searchText := "true"
	replacementText := "false"
	// get the command
	command := `go mod tidy &&
  golangci-lint run -c ./dev/golangci.toml --fix 2> /dev/null &&
  go test -rapid.nofailfile -failfast &&
  ./fuzz.fish`
	// Search the go files for mutation patterns
	files, _ := sh.Output("ag", searchText, "-G", `.*\.go$`, "-l")
	caught := true
	// For each file found
	for _, file := range strings.Split(files, "\n") {
		fmt.Printf("mutatable '%s' found in %s\n", searchText, file)
		candidates, _ := sh.Output("ag", "--column", searchText, file)
		// For each specific candidate found
		for _, candidate := range strings.Split(candidates, "\n") {
			//   replace the pattern
			numParts := 3
			parts := strings.SplitN(candidate, ":", numParts)
			line := parts[0]
			column, _ := strconv.Atoi(parts[1])
			match := parts[2]
			column--
			regex, _ := regexp.Compile(fmt.Sprintf(`(.{%d})%s`, column, searchText))
			mutant := regex.ReplaceAllString(match, fmt.Sprintf("${1}%s", replacementText))
			_ = sh.RunV("fish", "-c", fmt.Sprintf(
				`echo mutating %s:%s:(math %d+1) "%s" '->' "%s"`,
				file, line, column, match, mutant,
			))
			_ = sh.Run("fish", "-c", fmt.Sprintf(
				`sed -i "" -E %s's/(.{%d})%s/\1%s/' %s`,
				line, column, searchText, replacementText, file,
			))
			//   retest
			err := sh.RunV("fish", "-c", command)
			//   mark pass/failed
			if err == nil {
				_ = sh.RunV("fish", "-c", `echo failed to catch the mutant`)
				caught = false
			} else {
				_ = sh.RunV("fish", "-c", `echo caught the mutant`)
			}
			//   restore the pattern
			_ = sh.RunV("fish", "-c", fmt.Sprintf(
				`echo restoring mutant %s:%s:(math %d+1) "%s" '->' "%s"`,
				file, line, column, mutant, match,
			))
			_ = sh.Run("fish", "-c", fmt.Sprintf(
				`sed -i "" -E %s's/(.{%d})%s/\1%s/' %s`,
				line, column, replacementText, searchText, file,
			))
			//   if failed, exit
			if !caught {
				return
			}
			//   continue
			continue
		}
	}
}
