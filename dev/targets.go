//go:build targ

package dev

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toejough/targ"
)

func init() {
	targ.Register(
		targ.Targ("npx vite").Name("dev"),
		targ.Targ("npx vitest run").Name("test"),
		targ.Targ("npx vue-tsc --noEmit && npx vite build").Name("build"),
		targ.Targ("npx vue-tsc --noEmit && npx vitest run").Name("check"),
		targ.Targ("rm -rf dist/ test-results/").Name("clean"),
		targ.Targ(issues).Description("List open issues"),
		targ.Targ(issueClose).Description("Close an issue and commit"),
		targ.Targ(issueArchive).Description("Delete a closed issue and commit"),
		targ.Targ(history).Description("List deleted docs from git history"),
		targ.Targ(historyShow).Description("Show a deleted file from git history"),
	)
}

// issues lists all open issues.
func issues() error {
	entries, err := os.ReadDir("docs/issues")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No issues directory.")
			return nil
		}
		return err
	}

	found := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join("docs/issues", e.Name())
		status, err := readIssueStatus(path)
		if err != nil {
			continue
		}
		if status == "open" {
			fmt.Printf("  %s\n", strings.TrimSuffix(e.Name(), ".md"))
			found++
		}
	}

	if found == 0 {
		fmt.Println("No open issues.")
	} else {
		fmt.Printf("\n%d open issue(s).\n", found)
	}
	return nil
}

type closeArgs struct {
	Number string `targ:"positional,placeholder=NUMBER,desc=Issue number (e.g. 025)"`
}

func issueClose(args closeArgs) error {
	path, err := findIssueFile(args.Number)
	if err != nil {
		return err
	}

	status, err := readIssueStatus(path)
	if err != nil {
		return err
	}
	if status == "closed" {
		fmt.Printf("Issue %s is already closed.\n", args.Number)
		return nil
	}

	if err := verifyStatusEntry(args.Number); err != nil {
		return err
	}

	if err := replaceInFile(path, "**Status:** open", "**Status:** closed"); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	if err := targ.Run("git", "add", path); err != nil {
		return err
	}
	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	return targ.Run("git", "commit", "-m", fmt.Sprintf("docs: close #%s (%s)", args.Number, slug))
}

type archiveArgs struct {
	Number string `targ:"positional,placeholder=NUMBER,desc=Issue number (e.g. 025)"`
}

func issueArchive(args archiveArgs) error {
	path, err := findIssueFile(args.Number)
	if err != nil {
		return err
	}

	status, err := readIssueStatus(path)
	if err != nil {
		return err
	}
	if status == "open" {
		return fmt.Errorf("issue %s is still open — close it first with: targ issue-close %s", args.Number, args.Number)
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}

	if err := targ.Run("git", "add", path); err != nil {
		return err
	}
	return targ.Run("git", "commit", "-m", fmt.Sprintf("docs: archive #%s (%s)", args.Number, slug))
}

func history() error {
	return targ.RunV("git", "log", "--diff-filter=D", "--name-only", "--pretty=format:%h %s", "--", "docs/")
}

type showArgs struct {
	Path string `targ:"positional,placeholder=PATH,desc=Path to deleted file (e.g. docs/issues/001-bootstrap.md)"`
}

func historyShow(args showArgs) error {
	commitHash, err := targ.Output("git", "log", "--diff-filter=D", "--format=%H", "-1", "--", args.Path)
	if err != nil {
		return fmt.Errorf("could not find deletion commit for %s: %w", args.Path, err)
	}
	commitHash = strings.TrimSpace(commitHash)
	if commitHash == "" {
		return fmt.Errorf("no deletion found for %s — file may still exist or was never tracked", args.Path)
	}

	return targ.RunV("git", "show", commitHash+"~1:"+args.Path)
}

// Helper Functions

func findIssueFile(number string) (string, error) {
	entries, err := os.ReadDir("docs/issues")
	if err != nil {
		return "", fmt.Errorf("cannot read docs/issues/: %w", err)
	}
	prefix := number + "-"
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".md") {
			return filepath.Join("docs/issues", e.Name()), nil
		}
	}
	return "", fmt.Errorf("no issue file found matching %s-*.md in docs/issues/", number)
}

func readIssueStatus(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "**Status:**") {
			return strings.TrimSpace(strings.TrimPrefix(line, "**Status:**")), nil
		}
	}
	return "", fmt.Errorf("no **Status:** line found in %s", path)
}

func replaceInFile(path, old, new string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), old) {
		return fmt.Errorf("%q not found in %s", old, path)
	}
	updated := strings.Replace(string(data), old, new, 1)
	return os.WriteFile(path, []byte(updated), 0o644)
}

func verifyStatusEntry(number string) error {
	data, err := os.ReadFile("docs/status.md")
	if err != nil {
		return fmt.Errorf("cannot read docs/status.md: %w", err)
	}
	content := string(data)
	// Check both padded (#025) and unpadded (#25) forms
	unpadded := strings.TrimLeft(number, "0")
	if !strings.Contains(content, "#"+number) && !strings.Contains(content, "#"+unpadded) {
		return fmt.Errorf("docs/status.md has no reference to #%s — add a status entry before closing", number)
	}
	return nil
}
