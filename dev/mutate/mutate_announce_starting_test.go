package main

import (
	"spacer/dev/protest"
	"testing"
)

type (
	announceStartingDepsMock struct {
		calls *protest.FIFO[any]
		t     tester
	}
	printArgs struct{ message string }
	printCall struct{ args printArgs }
)

func newAnnounceStartingDepsMock(t tester) *announceStartingDepsMock {
	return &announceStartingDepsMock{
		calls: protest.NewFIFO[any]("calls"),
		t:     t,
	}
}

func (m *announceStartingDepsMock) print(s string) {
	m.calls.Push(printCall{args: printArgs{message: s}})
}

func (m *announceStartingDepsMock) close() {
	m.calls.Close()
}

func TestAnnounceStartingHappyPath(t *testing.T) {
	t.Parallel()

	deps := newAnnounceStartingDepsMock(t)

	// When the func is run
	go func() {
		(&runDepsMain{announceStartingDeps: deps}).announceStarting()
		deps.close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, printCall{args: printArgs{message: "Starting Mutation Testing"}})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
