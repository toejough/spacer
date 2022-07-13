package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"
)

type eqBool bool

func (e eqBool) Equals(b eqBool) bool {
	return e == b
}

func TestRun(t *testing.T) {
	t.Parallel()

	// linter can't tell we're using the range value in the test
	for _, mutationReturn := range []bool{true, false} { //nolint: paralleltest
		// Given a return from the mutation func
		mReturn := mutationReturn
		t.Run(fmt.Sprintf("%t", mutationReturn), func(t *testing.T) {
			t.Parallel()

			fut := protest.NewFUT()

			// Given dependencies
			mutateMock := protest.NewMockNoArgs[eqBool](fut, "mutate")
			mutate := func() bool {
				return bool(mutateMock.Func())
			}
			reportMock := protest.NewMockNoReturn[eqBool](fut, "report")
			report := func(r bool) {
				reportMock.Func(eqBool(r))
			}
			exitMock := protest.NewMockNoReturn[eqBool](fut, "exit")
			exit := func(r bool) {
				exitMock.Func(eqBool(r))
			}

			// When run is called
			fut.Call(func() {
				run(mutate, report, exit)
			})

			// Then we run the mutations
			mutateMock.WaitForCallFatal(t)

			// When we return the results of the mutation
			mutateMock.Return(eqBool(mReturn))

			// Then we report the summary of the run with the mutation results
			reportMock.ExpectCallFatal(t, eqBool(mReturn))

			// Then we exit with the result of the mutations
			exitMock.ExpectCallFatal(t, eqBool(mReturn))

			// Then we expect run to be done
			fut.ExpectDoneFatal(t)
		})
	}
}
