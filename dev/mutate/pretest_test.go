package main

import (
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

type pretestDepsMock struct {
	deps  pretestDeps
	calls *protest.FIFO[any]
}

type (
	fetchTestCommandCall protest.CallWithNoArgs[command]
	runTestCommandCall   protest.Call[command, bool]
)

func TestPretestHappyPath(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given test setup
		result, deps, fetchTestCommand := pretestTestSetup(test)

		// When the test command is returned
		testCommand := command(rapid.StringN(1, -1, -1).Draw(test, "test command"))
		fetchTestCommand.ReturnOneShot.Push(testCommand)

		// Then the test command is run
		var runTestCommand runTestCommandCall

		deps.calls.MustPopAs(test, &runTestCommand)
		protest.MustEqual(test, runTestCommand.Args, testCommand)

		// When the test command returns passing
		runTestCommand.ReturnOneShot.Push(true)

		// Then there are no more calls
		deps.calls.MustConfirmClosed(test)
		// And the function returns passing
		protest.MustEqual(test, true, *result)
	})
}

func TestPretestFetchCommandError(t *testing.T) {
	t.Parallel()

	// Given test setup
	result, deps, fetchTestCommand := pretestTestSetup(t)

	// When no command can be fetched
	fetchTestCommand.ReturnOneShot.Push("")

	// Then there are no more calls
	deps.calls.MustConfirmClosed(t)
	// And the function returns failing
	protest.MustEqual(t, false, *result)
}

func TestPretestCommandFailure(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given test setup
		result, deps, fetchTestCommand := pretestTestSetup(test)

		// When the test command is returned
		testCommand := command(rapid.StringN(1, -1, -1).Draw(test, "test command"))
		fetchTestCommand.ReturnOneShot.Push(testCommand)

		// Then the test command is run
		var runTestCommand runTestCommandCall

		deps.calls.MustPopAs(test, &runTestCommand)
		protest.MustEqual(test, runTestCommand.Args, testCommand)

		// When the test command returns failing
		runTestCommand.ReturnOneShot.Push(false)

		// Then there are no more calls
		deps.calls.MustConfirmClosed(test)
		// And the function returns failing
		protest.MustEqual(test, false, *result)
	})
}

func newPretestDepsMock(test tester) *pretestDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &pretestDepsMock{
		calls: calls,
		deps: pretestDeps{
			fetchTestCommand: func() command {
				return protest.ManageCallWithNoArgs[fetchTestCommandCall](test, calls)
			},
			runTestCommand: func(c command) bool { return protest.ManageCall[runTestCommandCall](test, calls, c) },
		},
	}
}

func pretestTestSetup(test protest.Tester) (*bool, *pretestDepsMock, fetchTestCommandCall) {
	// Given inputs/outputs
	var result bool

	deps := newPretestDepsMock(test)

	// When the function is called
	go func() {
		result = pretest(&deps.deps)
		deps.calls.Close()
	}()

	// Then the test command is fetched
	var fetchTestCommand fetchTestCommandCall

	deps.calls.MustPopAs(test, &fetchTestCommand)

	return &result, deps, fetchTestCommand
}
