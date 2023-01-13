package main

import (
	"spacer/dev/protest"
	"testing"
)

type (
	announcePretestDepsMock struct {
		calls *protest.FIFO[any]
		t     tester
		deps  announcePretestDeps
	}
)

func newAnnouncePretestDepsMock(t tester) *announcePretestDepsMock {
	calls := protest.NewFIFO[any]("calls")

	return &announcePretestDepsMock{
		calls: calls,
		t:     t,
		deps: announcePretestDeps{
			assumePrint: func(s string) { protest.ManageCallWithNoReturn[printCall](calls, s) },
		},
	}
}

func TestAnnouncePretestHappyPath(t *testing.T) {
	t.Parallel()

	deps := newAnnouncePretestDepsMock(t)

	// When the func is run
	go func() {
		announcePretest(&deps.deps)
		deps.calls.Close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, printCall{Args: "Starting Pretesting"})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
