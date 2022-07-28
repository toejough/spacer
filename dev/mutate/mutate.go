// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run{testFiles{newFileIterator, testAllPatterns}.f, reportResults, exit}.f()
}

type run struct {
	testFiles     func() bool
	reportResults func(bool)
	exit          func(int)
}

func (r run) f() {
	results := r.testFiles()
	r.reportResults(results)
	exitCode := resultsToExitCode(results)
	r.exit(exitCode)
}

type fileIterator interface {
	Next() string
}

type realFileIterator struct {
	filepaths []string
	i         int
}

func (r *realFileIterator) Next() string {
	if r.i >= len(r.filepaths) {
		return ""
	}

	current := r.filepaths[r.i]
	r.i++

	return current
}

func newFileIterator() fileIterator {
	return &realFileIterator{[]string{}, 0}
}

func resultsToExitCode(results bool) int {
	if !results {
		return 1
	}

	return 0
}

type testFiles struct {
	newFileIterator func() fileIterator
	testAllPatterns func(string) bool
}

func (tf testFiles) f() bool {
	iterator := tf.newFileIterator()
	path := iterator.Next()

	return tf.testAllPatterns(path)
}
func testAllPatterns(string) bool { return true }
func reportResults(bool)          {}
func exit(int)                    {}

//    searchText := "true"
//    replacementText := "false"
//    // get the command
//    commandPtr := flag.String("command", "go test ./...", "the command to run to attempt to catch the mutants.")
//    flag.Parse()

//    command := *commandPtr
//    // Search the go files for mutation patterns
//    matches, err := searchFiles(searchText)
//    if err != nil {
//        panic(fmt.Errorf("unable to mutate: %w", err))
//    }

//    caught := true
//    // For each file found
//    for _, match := range matches {
//        fmt.Printf("mutatable '%s' found at %d:%d(%s)\n", searchText, match.line, match.column, match.path)
//        //   replace the pattern
//        line := match.line
//        column := match.column
//        path := match.path
//        column--
//        _ = replaceText(line, column, searchText, replacementText, path)
//        //   retest
//        cp := strings.Fields(command)

//        var out []byte
//        out, err = exec.Command(cp[0], cp[1:]...).CombinedOutput() //nolint:gosec // I know I'm running user input.
//        fmt.Println(string(out))
//        //   mark pass/failed
//        if err == nil {
//            fmt.Printf("failed to catch the mutant\n")

//            caught = false
//        } else {
//            fmt.Printf("caught the mutant\n")
//        }
//        //   restore the pattern
//        _ = replaceText(line, column, replacementText, searchText, path)
//        //   if failed, exit
//        if !caught {
//            return false
//        }
//        //   continue
//        continue
//    }

//    return true
//}

// func report(result bool) {
//    if !result {
//        fmt.Println("Exiting early due to uncaught mutant")
//    } else {
//        fmt.Println("All mutants caught!")
//    }
//}

// func exit(result bool) {
//    if !result {
//        os.Exit(1)
//    } else {
//        os.Exit(0)
//    }
//}

//    input, err := ioutil.ReadFile(file)
//    if err != nil {
//        return fmt.Errorf("unable to replace text: unable to read file: %w", err)
//    }

//    lines := strings.Split(string(input), "\n")

//    targetLine := lines[line]
//    targetPart := targetLine[column:]
//    targetPart = strings.Replace(targetPart, searchText, replacementText, 1)
//    targetLine = targetLine[:column] + targetPart
//    lines[line] = targetLine
//    output := strings.Join(lines, "\n")

//    ownerReadWrite := 0o600

//    err = os.WriteFile(file, []byte(output), fs.FileMode(ownerReadWrite))
//    if err != nil {
//        return fmt.Errorf("unable to replace text: unable to write file: %w", err)
//    }

//    return nil
//}

// type match struct {
//    line   int
//    column int
//    path   string
//}

//    input, err := ioutil.ReadFile(file)
//    if err != nil {
//        return nil, fmt.Errorf("unable to replace text: unable to read file: %w", err)
//    }

//    lines := strings.Split(string(input), "\n")
//    matches := []match{}

//    for i, l := range lines {
//        indices := regexp.MustCompile(searchText).FindAllStringIndex(l, -1)
//        for _, pair := range indices {
//            matches = append(matches, match{line: i, column: pair[0], path: file})
//        }
//    }

//    return matches, nil
//}
//    workingDirectory, err := os.Getwd()
//    if err != nil {
//        return nil, fmt.Errorf("unable to search files: %w", err)
//    }

//    matches := []match{}

//    err = fs.WalkDir(os.DirFS(workingDirectory), ".", func(path string, _ fs.DirEntry, err error) error {
//        if err != nil {
//            fmt.Printf("error: %s\n", err)
//            return fs.SkipDir
//        }
//        if filepath.Ext(path) == ".go" {
//            m, _ := searchFile(searchText, path)
//            matches = append(matches, m...)
//        }

//        return nil
//    })

//    if err != nil {
//        return nil, fmt.Errorf("unable to search files: %w", err)
//    }

//    return matches, nil
//}
