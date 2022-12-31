package main

import (
	"spacer/dev/protest"
	"testing"
)

type verifyTestsPassWithNoMutantsDepsMock struct {
	deps  verifyTestsPassWithNoMutantsDeps
	calls *protest.FIFO[any]
}

type fetchTestCommandCall struct {
	returnOneShot *protest.FIFO[command]
}

type runTestCommandCall struct {
	args          command
	returnOneShot *protest.FIFO[bool]
}

func newVerifyTestsPassWithNoMutantsDepsMock(test tester) *verifyTestsPassWithNoMutantsDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &verifyTestsPassWithNoMutantsDepsMock{
		calls: calls,
		deps: verifyTestsPassWithNoMutantsDeps{
			fetchTestCommand: func() command {
				returnOneShot := protest.NewOneShotFIFO[command]("fetchTestCommand return")

				calls.Push(fetchTestCommandCall{
					returnOneShot: returnOneShot,
				})

				return returnOneShot.MustPop(test)
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

// TODO: stop using this function.
func (m *verifyTestsPassWithNoMutantsDepsMock) close() {
	m.calls.Close()
}

func TestVerifyTestsPassWithNoMutantsHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs/outputs...
	var result bool

	deps := newVerifyTestsPassWithNoMutantsDepsMock(t)

	// When the function is called
	go func() {
		result = verifyTestsPassWithNoMutants(&deps.deps)
		deps.close()
	}()

	// Then the test command is fetched
	var fetchTestCommand fetchTestCommandCall

	deps.calls.MustPopAs(t, &fetchTestCommand)

	// When the test command is returned
	// TODO rapid test this.
	fetchTestCommand.returnOneShot.Push("some command here")

	// Then the test command is run
	var runTestCommand runTestCommandCall

	deps.calls.MustPopAs(t, &runTestCommand)
	protest.MustEqual(t, runTestCommand.args, "some command here")

	// When the test command returns passing
	runTestCommand.returnOneShot.Push(true)

	// Then there are no more calls
	deps.calls.MustConfirmClosed(t)
	// And the function returns passing
	protest.MustEqual(t, true, result)
}
