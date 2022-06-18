package spacer_test

import (
	"spacer"
	"testing"
)

func TestAddGood(t *testing.T) {
	// Allow parallel testing
	t.Parallel()

	// GIVEN the command and good input
	goodInput := "add ABCD[E]FG"

	// WHEN the program is called with it
	ok := spacer.Run(goodInput)

	// THEN the result is ok
	if !ok {
		t.Fatal("Expected an ok result, but got not ok result")
	}
}

func FuzzAddGood(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		// WHEN the program is called with the input
		ok := spacer.Run(s)

		// THEN the result is expected to be ok.
		if !ok {
			t.Fatal("Expected an ok result, but got not ok result")
		}
	})
}
