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
	announcePretestCall  protest.CallWithNoArgsNoReturn
	fetchTestCommandCall protest.CallWithNoArgs[protest.Tuple[command]]
	runTestCommandCall   protest.Call[command, bool]
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
			runTestCommand: func(c command) bool { return protest.ManageCall[runTestCommandCall](test, calls, c) },
		},
	}
}

// TODO announcements from this function? stop/error?
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

		// Then there are no more calls
		deps.calls.MustConfirmClosed(test)
		// And the function returns passing
		protest.MustEqual(test, true, result)
	})
}

func TestVerifyTestsPassWithNoMutantsFetchCommandError(t *testing.T) {
	t.Parallel()

	// Given inputs/outputs
	var result bool

	deps := newPretestDepsMock(t)

	// When the function is called
	go func() {
		result = pretest(&deps.deps)
		deps.calls.Close()
	}()

	// Then the pretest is announced
	deps.calls.MustPopEqualTo(t, announcePretestCall{})
	// And the test command is fetched
	var fetchTestCommand fetchTestCommandCall

	deps.calls.MustPopAs(t, &fetchTestCommand)

	// When an error is returned
	// TODO rapid test the command & error
	fetchTestCommand.ReturnOneShot.Push(protest.Tuple[command]{
		Value: "arbitrary",
		// chill about dynamic error, this is a test
		Err: fmt.Errorf("arbitrary error"), //nolint: goerr113
	})

	// Then there are no more calls
	deps.calls.MustConfirmClosed(t)
	// And the function returns failing
	protest.MustEqual(t, false, result)
}

// TODO TestVerifyTestsPassWithNoMutantsRunCommandFailure
// TODO TestVerifyTestsPassWithNoMutantsRunCommandError
