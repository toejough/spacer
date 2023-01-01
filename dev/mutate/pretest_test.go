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

type fetchTestCommandCall protest.CallWithReturn[fetchTestCommandCallReturn]

type runTestCommandCall protest.CallWithArgsAndReturn[command, bool]

type fetchTestCommandCallReturn struct {
	command command
	err     error
}

type announcePretestCall protest.CallWithNoArgsNoReturn

func newPretestDepsMock(test tester) *pretestDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &pretestDepsMock{
		calls: calls,
		deps: pretestDeps{
			announcePretest: func() {
				calls.Push(announcePretestCall{})
			},
			fetchTestCommand: func() (command, error) {
				returnOneShot := protest.NewOneShotFIFO[fetchTestCommandCallReturn]("fetchTestCommand return")

				calls.Push(fetchTestCommandCall{
					ReturnOneShot: returnOneShot,
				})

				returnVals := returnOneShot.MustPop(test)

				return returnVals.command, returnVals.err
			},
			runTestCommand: func(c command) bool {
				returnOneShot := protest.NewOneShotFIFO[bool]("runTestCommand return")

				calls.Push(runTestCommandCall{
					Args:          c,
					ReturnOneShot: returnOneShot,
				})

				return returnOneShot.MustPop(test)
			},
		},
	}
}

// TODO some refactoring for how we set up these tests - they're still too tedius.
// * command types (no args/returns; args; returns; args & returns)
// * runner goroutine?
// * easier way to do the "as"
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
		fetchTestCommand.ReturnOneShot.Push(fetchTestCommandCallReturn{command: testCommand, err: nil})

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
	fetchTestCommand.ReturnOneShot.Push(fetchTestCommandCallReturn{
		command: "arbitrary",
		// chill about dynamic error, this is a test
		err: fmt.Errorf("arbitrary error"), //nolint: goerr113
	})

	// Then there are no more calls
	deps.calls.MustConfirmClosed(t)
	// And the function returns failing
	protest.MustEqual(t, false, result)
}

// TODO TestVerifyTestsPassWithNoMutantsRunCommandFailure
// TODO TestVerifyTestsPassWithNoMutantsRunCommandError
