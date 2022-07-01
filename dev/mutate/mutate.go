// Package mutate provides mutation testing functionality.
package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
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
			line, _ := strconv.Atoi(parts[0])
			column, _ := strconv.Atoi(parts[1])
			match := parts[2]
			column--
			regex, _ := regexp.Compile(fmt.Sprintf(`(.{%d})%s`, column, searchText))
			mutant := regex.ReplaceAllString(match, fmt.Sprintf("${1}%s", replacementText))

			fmt.Printf("mutating %s:%d:%d '%s' -> '%s'\n", file, line, column+1, match, mutant)

			_ = replaceText(line, column, searchText, replacementText, file)
			//   retest
			err := sh.RunV("fish", "-c", command)
			//   mark pass/failed
			if err == nil {
				fmt.Printf("failed to catch the mutant\n")

				caught = false
			} else {
				fmt.Printf("caught the mutant\n")
			}
			//   restore the pattern
			fmt.Printf("restoring mutant %s:%d:%d '%s' -> '%s'\n", file, line, column+1, mutant, match)

			_ = replaceText(line, column, replacementText, searchText, file)
			//   if failed, exit
			if !caught {
				return
			}
			//   continue
			continue
		}
	}
}

func replaceText(line int, column int, searchText string, replacementText string, file string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to replace text: unable to read file: %w", err)
	}

	lines := strings.Split(string(input), "\n")

	targetLine := lines[line-1]
	targetPart := targetLine[column:]
	targetPart = strings.Replace(targetPart, searchText, replacementText, 1)
	targetLine = targetLine[:column] + targetPart
	lines[line-1] = targetLine
	output := strings.Join(lines, "\n")

	ownerReadWrite := 0o600

	err = os.WriteFile(file, []byte(output), fs.FileMode(ownerReadWrite))
	if err != nil {
		return fmt.Errorf("unable to replace text: unable to write file: %w", err)
	}

	return nil
}
