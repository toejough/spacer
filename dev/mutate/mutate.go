// Package mutate provides mutation testing functionality.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

type (
	MutateFunc    func() bool
	ReportingFunc func(bool)
	ExitFunc      func(bool)
)

func run(m MutateFunc, r ReportingFunc, e ExitFunc) {
	mr := m()
	r(mr)
	e(mr)
}

func main() {
	searchText := "true"
	replacementText := "false"
	// get the command
	commandPtr := flag.String("command", "go test ./...", "the command to run to attempt to catch the mutants.")
	flag.Parse()

	command := *commandPtr
	// Search the go files for mutation patterns
	matches, err := searchFiles(searchText)
	if err != nil {
		panic(fmt.Errorf("unable to mutate: %w", err))
	}

	caught := true
	// For each file found
	for _, match := range matches {
		fmt.Printf("mutatable '%s' found at %d:%d(%s)\n", searchText, match.line, match.column, match.path)
		//   replace the pattern
		line := match.line
		column := match.column
		path := match.path
		column--
		_ = replaceText(line, column, searchText, replacementText, path)
		//   retest
		cp := strings.Fields(command)
		out, err := exec.Command(cp[0], cp[1:]...).CombinedOutput() //nolint:gosec // I know I'm running user input.
		fmt.Println(string(out))
		//   mark pass/failed
		if err == nil {
			fmt.Printf("failed to catch the mutant\n")

			caught = false
		} else {
			fmt.Printf("caught the mutant\n")
		}
		//   restore the pattern
		_ = replaceText(line, column, replacementText, searchText, path)
		//   if failed, exit
		if !caught {
			fmt.Println("Exiting early due to uncaught mutant")

			return
		}
		//   continue
		continue
	}

	fmt.Println("All mutants caught!")
}

func replaceText(line int, column int, searchText string, replacementText string, file string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to replace text: unable to read file: %w", err)
	}

	lines := strings.Split(string(input), "\n")

	targetLine := lines[line]
	targetPart := targetLine[column:]
	targetPart = strings.Replace(targetPart, searchText, replacementText, 1)
	targetLine = targetLine[:column] + targetPart
	lines[line] = targetLine
	output := strings.Join(lines, "\n")

	ownerReadWrite := 0o600

	err = os.WriteFile(file, []byte(output), fs.FileMode(ownerReadWrite))
	if err != nil {
		return fmt.Errorf("unable to replace text: unable to write file: %w", err)
	}

	return nil
}

type match struct {
	line   int
	column int
	path   string
}

func searchFile(searchText string, file string) ([]match, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to replace text: unable to read file: %w", err)
	}

	lines := strings.Split(string(input), "\n")
	matches := []match{}

	for i, l := range lines {
		indices := regexp.MustCompile(searchText).FindAllStringIndex(l, -1)
		for _, pair := range indices {
			matches = append(matches, match{line: i, column: pair[0], path: file})
		}
	}

	return matches, nil
}

func searchFiles(searchText string) ([]match, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to search files: %w", err)
	}

	matches := []match{}

	err = fs.WalkDir(os.DirFS(workingDirectory), ".", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("error: %s\n", err)

			return fs.SkipDir
		}
		if filepath.Ext(path) == ".go" {
			m, _ := searchFile(searchText, path)
			matches = append(matches, m...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to search files: %w", err)
	}

	return matches, nil
}
