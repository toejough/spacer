//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	//"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// better glob expansion
// https://stackoverflow.com/a/26809999
func globs(dir string, ext []string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to find all glob matches: %w", err)
		}
		for _, each := range ext {
			if filepath.Ext(path) == each {
				files = append(files, path)
                return nil
			}
		}
		return nil
	})

	return files, err
}

// Run all checks on the code whenever a relevant file changes independently
func Monitor() error {
	fmt.Println("Monitoring...")

	err := Check()
	if err != nil {
		fmt.Printf("continuing to monitor after check failure: %s", err)
	}

	lastFinishedTime := time.Now()

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to monitor effectively due to error getting current working directory: %w", err)
	}

	for {
		time.Sleep(time.Second)

		paths, err := globs(dir, []string{".go", ".fish", ".toml"})
		if err != nil {
			return fmt.Errorf("unable to monitor effectively due to error resolving globs: %w", err)
		}

		changeDetected, err := target.PathNewer(lastFinishedTime, paths...)
		if err != nil {
			return fmt.Errorf("unable to monitor effectively due to error checking for path updates: %w", err)
		}

		if changeDetected {
			fmt.Println("Change detected...")
			err = Check()
            if err != nil {
                fmt.Printf("continuing to monitor after check failure: %s", err)
            }

			lastFinishedTime = time.Now()
		}
	}
}

// Run all checks on the code
func Check() error {
	fmt.Println("Checking...")
	for _, cmd := range []func() error{Tidy, Lint, Test, Fuzz, Mutate} {
		err := cmd()
		if err != nil {
			return fmt.Errorf("unable to finish checking: %w", err)
		}
	}
	return nil
}

// Run all checks on the code for determining whether any fail
func CheckForFail() error {
	fmt.Println("Checking...")
	for _, cmd := range []func() error{LintForFail, TestForFail, Fuzz} {
		err := cmd()
		if err != nil {
			return fmt.Errorf("unable to finish checking: %w", err)
		}
	}
	return nil
}

// Tidy tidies up go.mod
func Tidy() error {
	fmt.Println("Tidying go.mod...")
	return sh.RunV("go", "mod", "tidy")
}

// Lint lints the codebase
func Lint() error {
	fmt.Println("Linting...")
	_, err := sh.Exec(nil, os.Stdout, nil, "golangci-lint", "run", "-c", "dev/golangci.toml")
	return err
}

// LintForFail lints the codebase purely to find out whether anything fails
func LintForFail() error {
	fmt.Println("Linting to check for overall pass/fail...")
	_, err := sh.Exec(
        nil, os.Stdout, nil,
        "golangci-lint", "run",
        "-c", "dev/golangci.toml",
        "--fix=false",
        "--max-issues-per-linter=1",
        "--max-same-issues=1",
    )
	return err
}

// Run the unit tests
func Test() error {
	fmt.Println("Running unit tests...")
	return sh.RunV("go", "test", "./...")
}

// Run the unit tests purely to find out whether any fail
func TestForFail() error {
	fmt.Println("Running unit tests for overall pass/fail...")
	return sh.RunV("go", "test", "./...", "-rapid.nofailfile", "-failfast")
}

// Run the fuzz tests
func Fuzz() error {
	fmt.Println("Running fuzz tests...")
	return sh.RunV("./dev/fuzz.fish")
}

// Run the mutation tests
func Mutate() error {
	fmt.Println("Running mutation tests...")
	return sh.RunV("go", "run", "./dev/mutate/mutate.go", "--command", "mage checkForFail")
}

// Install development tooling
func InstallTools() error {
	fmt.Println("Installing development tools...")
	return sh.RunV("./dev/dev-install.sh")
}

// Clean up the dev env
func Clean() {
	fmt.Println("Cleaning...")
	os.Remove("coverage.out")
}
