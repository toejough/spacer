package main

import (
	"testing"
	"time"
)

type MockNoArgs struct {
	t          *testing.T
	callsChan  chan string
	resultChan chan bool
	name       string
}

func (f *FUT) NewMockNoArgs(name string) *MockNoArgs {
	f.T.Helper()

	mock := &MockNoArgs{
		t:          f.T,
		callsChan:  f.CallsChan,
		resultChan: make(chan bool),
		name:       name,
	}

	return mock
}

func (m *MockNoArgs) ExpectCall() *MockNoArgs {
	m.t.Helper()

	waitForCall(m.t, m.callsChan, m.name)

	return m
}

func (m *MockNoArgs) Return(r bool) *MockNoArgs {
	m.t.Helper()

	setReturnValue(m.t, m.resultChan, r)

	return m
}

func (m *MockNoArgs) Func() bool {
	m.callsChan <- m.name

	return <-m.resultChan
}

type MockNoReturn struct {
	t         *testing.T
	callsChan chan string
	argsChan  chan bool
	name      string
}

func (f *FUT) NewMockNoReturn(name string) *MockNoReturn {
	f.T.Helper()

	mock := &MockNoReturn{
		t:         f.T,
		callsChan: f.CallsChan,
		argsChan:  make(chan bool),
		name:      name,
	}

	return mock
}

func (m *MockNoReturn) ExpectCall(b bool) *MockNoReturn {
	m.t.Helper()

	waitForCall(m.t, m.callsChan, m.name)
	waitForArgs(m.t, m.argsChan, b, m.name)

	return m
}

func (m *MockNoReturn) Func(args bool) {
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

func TestRun(t *testing.T) {
	t.Parallel()

	fut := NewFUT(t)

	// Given a mutation func
	mutateMock := fut.NewMockNoArgs("mutate")
	mutate := func() bool {
		return mutateMock.Func()
	}

	// Given a reporting func
	reportMock := fut.NewMockNoReturn("report")
	report := func(r bool) {
		reportMock.Func(r)
	}

	// Given an exit func
	exitMock := fut.NewMockNoReturn("exit")
	exit := func(r bool) {
		exitMock.Func(r)
	}

	// Given a return from the mutation func
	mutationReturn := true

	// When run is called
	fut.Call(func() {
		run(mutate, report, exit)
	})

	// Then we run the mutations and return the results
	mutateMock.ExpectCall().Return(mutationReturn)

	// Then we report the summary of the run with the mutation results
	reportMock.ExpectCall(mutationReturn)

	// Then we exit with the result of the mutations
	exitMock.ExpectCall(mutationReturn)

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
