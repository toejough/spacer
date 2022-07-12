package main

import (
	"fmt"
	"testing"
	"time"
)

type Eq[E any] interface {
	Equals(E) bool
}

type MockNoArgs struct {
	callsChan  chan string
	resultChan chan bool
	name       string
}

func (f *FUT) NewMockNoArgs(name string) *MockNoArgs {
	return &MockNoArgs{
		callsChan:  f.CallsChan,
		resultChan: make(chan bool),
		name:       name,
	}
}

func (m *MockNoArgs) WaitForCall() error {
	return waitForCall(m.callsChan, m.name)
}

func (m *MockNoArgs) WaitForCallFatal(t *testing.T) {
	t.Helper()

	err := m.WaitForCall()
	if err != nil {
		t.Fatalf("did not call '%s': %s", m.name, err)
	}
}

func (m *MockNoArgs) Return(r bool) {
	m.resultChan <- r
}

func (m *MockNoArgs) Func() bool {
	m.callsChan <- m.name

	return <-m.resultChan
}

type MockNoReturn[A Eq[A]] struct {
	callsChan chan string
	argsChan  chan A
	name      string
}

func NewMockNoReturn[A Eq[A]](f *FUT, name string) *MockNoReturn[A] {
	return &MockNoReturn[A]{
		callsChan: f.CallsChan,
		argsChan:  make(chan A),
		name:      name,
	}
}

func (mock *MockNoReturn[A]) ExpectCall() (*A, error) {
	err := waitForCall(mock.callsChan, mock.name)
	if err != nil {
		return nil, err
	}

	return waitForArgs(mock.argsChan, mock.name)
}

func (mock *MockNoReturn[A]) ExpectCallFatal(t *testing.T, expected A) {
	t.Helper()

	args, err := mock.ExpectCall()
	if err != nil {
		t.Fatalf("'%s' was not called with the expected args: %s", mock.name, err)
	}

	if args == nil {
		t.Fatalf("args were unexpectedly nil. expected: %v", expected)
	}

	if !(*args).Equals(expected) {
		t.Fatalf(
			"Expected 'report' to be called with arguments of '%v', but got arguments of '%v' instead\n",
			expected, *args,
		)
	}
}

func (mock *MockNoReturn[A]) Func(args A) {
	mock.callsChan <- mock.name
	mock.argsChan <- args
}

type FUT struct {
	CallsChan chan string
}

func NewFUT() *FUT {
	return &FUT{CallsChan: make(chan string)}
}

type UnexpectedCallError struct {
	Call string
}

func (e UnexpectedCallError) Error() string {
	return fmt.Sprintf("expected to be done, but function called '%s' instead", e.Call)
}

type UnendingError struct{}

func (e UnendingError) Error() string {
	return "expected to be done, but function timed out instead"
}

func (f *FUT) ExpectDone() error {
	select {
	case call := <-f.CallsChan:
		if call != "" {
			return UnexpectedCallError{call}
		}

		return nil
	case <-time.After(time.Second):
		return UnendingError{}
	}
}

func (f *FUT) ExpectDoneFatal(t *testing.T) {
	t.Helper()

	err := f.ExpectDone()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func (f *FUT) Call(ff func()) {
	go func() {
		ff()
		close(f.CallsChan)
	}()
}

type eqBool bool

func (e eqBool) Equals(b eqBool) bool {
	return e == b
}

func TestRun(t *testing.T) {
	t.Parallel()

	fut := NewFUT()

	// Given dependencies
	mutateMock := fut.NewMockNoArgs("mutate")
	mutate := func() bool {
		return mutateMock.Func()
	}
	reportMock := NewMockNoReturn[eqBool](fut, "report")
	report := func(r bool) {
		reportMock.Func(eqBool(r))
	}
	exitMock := NewMockNoReturn[eqBool](fut, "exit")
	exit := func(r bool) {
		exitMock.Func(eqBool(r))
	}

	// Given a return from the mutation func
	mutationReturn := true

	// When run is called
	fut.Call(func() {
		run(mutate, report, exit)
	})

	// Then we run the mutations
	mutateMock.WaitForCallFatal(t)

	// When we return the results of the mutation
	mutateMock.Return(mutationReturn)

	// Then we report the summary of the run with the mutation results
	reportMock.ExpectCallFatal(t, eqBool(mutationReturn))

	// Then we exit with the result of the mutations
	exitMock.ExpectCallFatal(t, eqBool(mutationReturn))

	// Then we expect run to be done
	fut.ExpectDoneFatal(t)
}

type CallError struct {
	Expected, Actual string
}

func (e CallError) Error() string {
	return fmt.Sprintf("expected '%s' to be called, but '%s' was called instead", e.Expected, e.Actual)
}

type CallTimeoutError struct {
	Expected string
}

func (e CallTimeoutError) Error() string {
	return fmt.Sprintf("expected run to call '%s' before timing out, but it did not", e.Expected)
}

type ArgTimeoutError struct{}

func (e ArgTimeoutError) Error() string {
	return "expected to receive arguments before timing out, but did not"
}

func waitForCall(callsChan chan string, call string) error {
	select {
	case called := <-callsChan:
		if called != call {
			return CallError{call, called}
		}

		return nil
	case <-time.After(time.Second):
		return CallTimeoutError{call}
	}
}

func waitForArgs[E any](argsChan chan E, toFunc string) (*E, error) {
	select {
	case actualArgs := <-argsChan:
		return &actualArgs, nil
	case <-time.After(time.Second):
		return nil, ArgTimeoutError{}
	}
}
