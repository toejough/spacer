// Package protest provides procedure testing functionality.
package protest

import (
	"fmt"
	"testing"
	"time"
)

type Eq[E any] interface {
	Equals(E) bool
}

type MockNoArgs[A any] struct {
	callsChan  chan string
	resultChan chan A
	name       string
}

func NewMockNoArgs[A any](f *FUT, name string) *MockNoArgs[A] {
	return &MockNoArgs[A]{
		callsChan:  f.CallsChan,
		resultChan: make(chan A),
		name:       name,
	}
}

func (mock *MockNoArgs[A]) WaitForCall() error {
	return waitForCall(mock.callsChan, mock.name)
}

func (mock *MockNoArgs[A]) WaitForCallFatal(t *testing.T) {
	t.Helper()

	err := mock.WaitForCall()
	if err != nil {
		t.Fatalf("did not call '%s': %s", mock.name, err)
	}
}

func (mock *MockNoArgs[A]) Return(r A) {
	mock.resultChan <- r
}

// I know this returns an interface - it's a generic return type.
func (mock *MockNoArgs[A]) Func() A { //nolint: ireturn
	mock.callsChan <- mock.name

	return <-mock.resultChan
}

type MockNoReturn[A any] struct {
	callsChan chan string
	argsChan  chan A
	name      string
}

func NewMockNoReturn[A any](f *FUT, name string) *MockNoReturn[A] {
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

func (mock *MockNoReturn[A]) ExpectCallFatal(t *testing.T, expected Eq[A]) {
	t.Helper()

	args, err := mock.ExpectCall()
	if err != nil {
		t.Fatalf("'%s' was not called with the expected args: %s", mock.name, err)
	}

	if args == nil {
		t.Fatalf("args were unexpectedly nil. expected: %v", expected)
	}

	if !expected.Equals(*args) {
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
