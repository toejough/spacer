package main

import (
	"spacer/dev/protest"
	"testing"
)

type (
	announceStartingDepsMock struct {
		calls *protest.FIFO[any]
		t     tester
		deps  announceStartingDeps
	}
	printArgs struct{ message string }
	printCall protest.CallWithArgs[printArgs]
)

func newAnnounceStartingDepsMock(t tester) *announceStartingDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &announceStartingDepsMock{
		calls: calls,
		t:     t,
		deps: announceStartingDeps{
			print: func(s string) {
				calls.Push(printCall{Args: printArgs{message: s}})
			},
		},
	}
}

func (m *announceStartingDepsMock) close() {
	m.calls.Close()
}

func TestAnnounceStartingHappyPath(t *testing.T) {
	t.Parallel()

	deps := newAnnounceStartingDepsMock(t)

	// When the func is run
	go func() {
		announceStarting(&deps.deps)
		deps.close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, printCall{Args: printArgs{message: "Starting Mutation Testing"}})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
