package spacer_test

import (
	"regexp"
	"spacer/pkg/spacer"
	"testing"

	"pgregory.net/rapid"
)

func TestAddGood(t *testing.T) {
	// Allow parallel testing
	t.Parallel()

	// GIVEN good input
	goodInput := "add ABCD[E]FG"

	// WHEN the program is called with it
	ok := spacer.Run(goodInput)

	// THEN the result is ok
	if !ok {
		t.Fatal("Expected an ok result, but got not ok result")
	}
}

func FuzzAdd(f *testing.F) {
	f.Add("ABCD[E]FG") // ok example
	f.Add("ABCDEFG")   // not ok example
	f.Fuzz(func(t *testing.T, s string) {
		// WHEN the program is called with the input
		_ = spacer.Run(s)

		// THEN the run is expected to return just fine.
	})
}

func TestPropertyAddGood(t *testing.T) {
	// Allow parallel testing
	t.Parallel()

	// Perform the property check
	rapid.Check(t, func(rtest *rapid.T) {
		// GIVEN good input
		gen := rapid.StringMatching(`.*(\[.+\].*)+`)
		goodInput := gen.Draw(rtest, "goodInput")

		// WHEN the program is called with it
		ok := spacer.Run(goodInput)

		// THEN the result is ok
		if !ok {
			rtest.Fatal("Expected an ok result, but got not ok result")
		}
	})
}

func TestAddBad(t *testing.T) {
	// Allow parallel testing
	t.Parallel()

	// GIVEN bad input
	badInput := "add ABCDEFG"

	// WHEN the program is called with it
	ok := spacer.Run(badInput)

	// THEN the result is not ok
	if ok {
		t.Fatal("Expected result to not be ok, but it was ok")
	}
}

func TestPropertyAddBad(t *testing.T) {
	// Allow parallel testing
	t.Parallel()

	regex := regexp.MustCompile(`.*\[.+\].*`)

	// Perform the property check
	rapid.Check(t, func(rtest *rapid.T) {
		// GIVEN bad input
		gen := rapid.String().Filter(func(s string) bool { return !regex.MatchString(s) })
		badInput := gen.Draw(rtest, "badInput")

		// WHEN the program is called with it
		ok := spacer.Run(badInput)

		// THEN the result is not ok
		if ok {
			rtest.Fatal("Expected result not to be ok, but it was ok")
		}
	})
}
