package main

import (
	"testing"
	"time"
)

type Mock struct {
	t          *testing.T
	callsChan  chan string
	resultChan chan bool
	name       string
}

func NewMock(t *testing.T, callsChan chan string, name string) *Mock {
	t.Helper()

	mock := &Mock{
		t:          t,
		callsChan:  callsChan,
		resultChan: make(chan bool),
		name:       name,
	}

	return mock
}

func (m *Mock) ExpectCall() *Mock {
	m.t.Helper()

	waitForCall(m.t, m.callsChan, m.name)

	return m
}

func (m *Mock) Return(r bool) *Mock {
	m.t.Helper()

	setReturnValue(m.t, m.resultChan, r)

	return m
}

func (m *Mock) Func() bool {
	m.callsChan <- m.name

	return <-m.resultChan
}

type MockArgs struct {
	t         *testing.T
	callsChan chan string
	argsChan  chan bool
	name      string
}

func NewMockArgs(t *testing.T, callsChan chan string, name string) *MockArgs {
	t.Helper()

	mock := &MockArgs{
		t:         t,
		callsChan: callsChan,
		argsChan:  make(chan bool),
		name:      name,
	}

	return mock
}

func (m *MockArgs) ExpectCall(b bool) *MockArgs {
	m.t.Helper()

	waitForCall(m.t, m.callsChan, m.name)
	waitForArgs(m.t, m.argsChan, b, m.name)

	return m
}

func (m *MockArgs) Func(args bool) {
	m.callsChan <- m.name
	m.argsChan <- args
}

type FUT struct {
	CallsChan chan string
	T         *testing.T
}

func NewFUT(t *testing.T) *FUT {
	t.Helper()

	return &FUT{T: t, CallsChan: make(chan string)}
}

func (f *FUT) ExpectDone() {
	f.T.Helper()
	select {
	case call := <-f.CallsChan:
		if call != "" {
			f.T.Fatalf("Expected run to be done, but it called '%s' instead.\n", call)
		}
	case <-time.After(time.Second):
		f.T.Fatalf("Expected run to be done, but it didn't end.\n")
	}
}

func (f *FUT) Call(ff func()) {
	go func() {
		ff()
		close(f.CallsChan)
	}()
}

func TestMainGetsCommand(t *testing.T) {
	t.Parallel()

	fut := NewFUT(t)

	// Given a mutation func
	// TODO make the new mocks functions of fut, to further hide the call channel
	mutateMock := NewMock(t, fut.CallsChan, "mutate")
	mutate := func() bool {
		return mutateMock.Func()
	}

	// Given a reporting func
	reportMock := NewMockArgs(t, fut.CallsChan, "report")
	report := func(r bool) {
		reportMock.Func(r)
	}

	// Given an exit func
	exitMock := NewMockArgs(t, fut.CallsChan, "exit")
	exit := func(r bool) {
		exitMock.Func(r)
	}

	// When run is called
	fut.Call(func() {
		run(mutate, report, exit)
	})

	// Then we run the mutations and return the results
	mutateMock.ExpectCall().Return(true)

	// Then we report the summary of the run with the mutation results
	reportMock.ExpectCall(true)

	// Then we exit with the result of the mutations
	exitMock.ExpectCall(true)

	// Then we expect run to be done
	fut.ExpectDone()
}

func waitForCall(t *testing.T, callsChan chan string, call string) {
	t.Helper()
	select {
	case called := <-callsChan:
		if called != call {
			t.Fatalf("Expected '%s' to be called, but '%s' was called instead\n", call, called)
		}
	case <-time.After(time.Second):
		t.Fatalf("Expected run to call '%s' before timing out, but it did not.\n", call)
	}
}

func waitForArgs(t *testing.T, argsChan chan bool, expectedArgs bool, toFunc string) {
	t.Helper()
	select {
	case actualArgs := <-argsChan:
		if actualArgs != expectedArgs {
			t.Fatalf(
				"Expected '%s' to be called with arguments of '%t', but got arguments of '%t' instead\n",
				toFunc, expectedArgs, actualArgs,
			)
		}
	case <-time.After(time.Second):
		t.Fatalf(
			"Expected '%s' to be called with arguments of '%t' before timing out, but it was not.\n",
			toFunc, expectedArgs,
		)
	}
}

func setReturnValue(t *testing.T, returnChan chan bool, returnValue bool) {
	t.Helper()
	returnChan <- returnValue
}
