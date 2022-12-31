package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

type verifyTestsPassWithNoMutantsDepsMock struct {
	deps  verifyTestsPassWithNoMutantsDeps
	calls *protest.FIFO[any]
}

type fetchTestCommandCall struct {
	returnOneShot *protest.FIFO[fetchTestCommandCallReturn]
}

type runTestCommandCall struct {
	args          command
	returnOneShot *protest.FIFO[bool]
}

type fetchTestCommandCallReturn struct {
	command command
	err     error
}

func newVerifyTestsPassWithNoMutantsDepsMock(test tester) *verifyTestsPassWithNoMutantsDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &verifyTestsPassWithNoMutantsDepsMock{
		calls: calls,
		deps: verifyTestsPassWithNoMutantsDeps{
			fetchTestCommand: func() (command, error) {
				returnOneShot := protest.NewOneShotFIFO[fetchTestCommandCallReturn]("fetchTestCommand return")

				calls.Push(fetchTestCommandCall{
					returnOneShot: returnOneShot,
				})

				returnVals := returnOneShot.MustPop(test)

				return returnVals.command, returnVals.err
			},
			runTestCommand: func(c command) bool {
				returnOneShot := protest.NewOneShotFIFO[bool]("runTestCommand return")

				calls.Push(runTestCommandCall{
					args:          c,
					returnOneShot: returnOneShot,
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
// TODO announcements from this function? start/stop/error?
func TestVerifyTestsPassWithNoMutantsHappyPath(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs...
		var result bool

		deps := newVerifyTestsPassWithNoMutantsDepsMock(test)

		// When the function is called
		go func() {
			result = verifyTestsPassWithNoMutants(&deps.deps)
			deps.calls.Close()
		}()

		// Then the test command is fetched
		var fetchTestCommand fetchTestCommandCall

		deps.calls.MustPopAs(test, &fetchTestCommand)

		// When the test command is returned
		testCommand := command(rapid.String().Draw(test, "test command"))
		fetchTestCommand.returnOneShot.Push(fetchTestCommandCallReturn{command: testCommand, err: nil})

		// Then the test command is run
		var runTestCommand runTestCommandCall

		deps.calls.MustPopAs(test, &runTestCommand)
		protest.MustEqual(test, runTestCommand.args, testCommand)

		// When the test command returns passing
		runTestCommand.returnOneShot.Push(true)

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

	deps := newVerifyTestsPassWithNoMutantsDepsMock(t)

	// When the function is called
	go func() {
		result = verifyTestsPassWithNoMutants(&deps.deps)
		deps.calls.Close()
	}()

	// Then the test command is fetched
	var fetchTestCommand fetchTestCommandCall

	deps.calls.MustPopAs(t, &fetchTestCommand)

	// When an error is returned
	// TODO rapid test the command & error
	fetchTestCommand.returnOneShot.Push(fetchTestCommandCallReturn{
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
