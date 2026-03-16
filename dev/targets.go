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
		targ.Targ(issueNew).Description("Scaffold a new issue file with next number"),
		targ.Targ(issueClose).Description("Close an issue, archive it, and commit"),
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

type newArgs struct {
	Slug string `targ:"positional,placeholder=SLUG,desc=Issue slug (e.g. fix-login-bug)"`
}

func issueNew(args newArgs) error {
	if args.Slug == "" {
		return fmt.Errorf("slug is required")
	}

	// Ensure docs/issues/ exists
	if err := os.MkdirAll("docs/issues", 0o755); err != nil {
		return fmt.Errorf("failed to create docs/issues/: %w", err)
	}

	// Find next issue number by scanning existing files and git history
	next, err := nextIssueNumber()
	if err != nil {
		return err
	}

	number := fmt.Sprintf("%03d", next)
	filename := fmt.Sprintf("%s-%s.md", number, args.Slug)
	path := filepath.Join("docs/issues", filename)

	template := fmt.Sprintf(`# %s — %s

**Status:** open
**Type:**

## Problem



## Principle



## Guidance

Before implementing, read the current codebase to understand what's changed since this issue was filed. Research external best practices relevant to the problem. Tailor the solution to the current state, not the state when this issue was written.
`, number, toTitle(args.Slug))

	if err := os.WriteFile(path, []byte(template), 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	fmt.Printf("Created %s\n", path)
	return nil
}

func nextIssueNumber() (int, error) {
	max := 0

	// Check current files
	entries, err := os.ReadDir("docs/issues")
	if err != nil && !os.IsNotExist(err) {
		return 0, fmt.Errorf("cannot read docs/issues/: %w", err)
	}
	for _, e := range entries {
		n := parseIssueNumber(e.Name())
		if n > max {
			max = n
		}
	}

	// Check git history for deleted issues
	output, err := targ.Output("git", "log", "--diff-filter=D", "--name-only", "--pretty=format:", "--", "docs/issues/")
	if err == nil {
		for _, line := range strings.Split(output, "\n") {
			base := filepath.Base(strings.TrimSpace(line))
			n := parseIssueNumber(base)
			if n > max {
				max = n
			}
		}
	}

	return max + 1, nil
}

func parseIssueNumber(filename string) int {
	parts := strings.SplitN(filename, "-", 2)
	if len(parts) < 2 {
		return 0
	}
	n := 0
	for _, c := range parts[0] {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func toTitle(slug string) string {
	words := strings.Split(slug, "-")
	if len(words) > 0 {
		words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	}
	return strings.Join(words, " ")
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

	// Gate 1: working tree must be clean (all related work committed)
	if err := verifyCleanWorkingTree(); err != nil {
		return err
	}

	// Gate 2: status.md must reference this issue
	if err := verifyStatusEntry(args.Number); err != nil {
		return err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")

	// Close: update status and commit
	if err := replaceInFile(path, "**Status:** open", "**Status:** closed"); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	if err := targ.Run("git", "add", path); err != nil {
		return err
	}
	if err := targ.Run("git", "commit", "-m", fmt.Sprintf("docs: close #%s (%s)", args.Number, slug)); err != nil {
		return err
	}

	// Archive: delete file and commit
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

func verifyCleanWorkingTree() error {
	output, err := targ.Output("git", "status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to check working tree: %w", err)
	}
	if strings.TrimSpace(output) != "" {
		return fmt.Errorf("working tree has uncommitted changes — commit or stash them first, then re-run:\n%s", output)
	}
	return nil
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
		return fmt.Errorf("#%s not found in status.md — add it to the current cycle's section, or start a new cycle if this is new work", number)
	}
	return nil
}
