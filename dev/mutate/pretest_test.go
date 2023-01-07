package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

type pretestDepsMock struct {
	deps  pretestDeps
	calls *protest.FIFO[any]
}

type (
	announcePretestCall        protest.CallWithNoArgsNoReturn
	fetchTestCommandCall       protest.CallWithNoArgs[protest.Tuple[command]]
	runTestCommandCall         protest.Call[command, bool]
	announcePretestResultsCall protest.CallWithNoReturn[bool]
)

func newPretestDepsMock(test tester) *pretestDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &pretestDepsMock{
		calls: calls,
		deps: pretestDeps{
			announcePretest: func() { protest.ManageCallWithNoArgsNoReturn[announcePretestCall](calls) },
			fetchTestCommand: func() (command, error) {
				return protest.ManageCallWithNoArgs[fetchTestCommandCall](test, calls).Unwrap() //nolint: wrapcheck
			},
			runTestCommand:         func(c command) bool { return protest.ManageCall[runTestCommandCall](test, calls, c) },
			announcePretestResults: func(b bool) { protest.ManageCallWithNoReturn[announcePretestResultsCall](calls, b) },
		},
	}
}

func TestVerifyTestsPassWithNoMutantsHappyPath(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs...
		var result bool

		deps := newPretestDepsMock(test)

		// When the function is called
		go func() {
			result = pretest(&deps.deps)
			deps.calls.Close()
		}()

		// Then the pretest is announced
		deps.calls.MustPopEqualTo(test, announcePretestCall{})
		// and the test command is fetched
		var fetchTestCommand fetchTestCommandCall

		deps.calls.MustPopAs(test, &fetchTestCommand)

		// When the test command is returned
		testCommand := command(rapid.String().Draw(test, "test command"))
		fetchTestCommand.ReturnOneShot.Push(protest.Tuple[command]{Value: testCommand, Err: nil})

		// Then the test command is run
		var runTestCommand runTestCommandCall

		deps.calls.MustPopAs(test, &runTestCommand)
		protest.MustEqual(test, runTestCommand.Args, testCommand)

		// When the test command returns passing
		runTestCommand.ReturnOneShot.Push(true)

		// Then the pretest result is announced
		deps.calls.MustPopEqualTo(test, announcePretestResultsCall{Args: true})

		// Then there are no more calls
		deps.calls.MustConfirmClosed(test)
		// And the function returns passing
		protest.MustEqual(test, true, result)
	})
}

func TestVerifyTestsPassWithNoMutantsFetchCommandError(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result bool

		deps := newPretestDepsMock(test)

		// When the function is called
		go func() {
			result = pretest(&deps.deps)
			deps.calls.Close()
		}()

		// Then the pretest is announced
		deps.calls.MustPopEqualTo(test, announcePretestCall{})
		// And the test command is fetched
		var fetchTestCommand fetchTestCommandCall

		deps.calls.MustPopAs(test, &fetchTestCommand)

		// When an error is returned
		fetchTestCommand.ReturnOneShot.Push(protest.Tuple[command]{
			Value: command(rapid.String().Draw(test, "test command")),
			// chill about dynamic error, this is a test
			Err: fmt.Errorf(rapid.String().Draw(test, "test error")), //nolint: goerr113
		})

		// Then there are no more calls
		deps.calls.MustConfirmClosed(test)
		// And the function returns failing
		protest.MustEqual(test, false, result)
	})
}

// TODO TestVerifyTestsPassWithNoMutantsRunCommandFailure
// TODO TestVerifyTestsPassWithNoMutantsRunCommandError
