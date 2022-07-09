package main

import (
	"testing"
	"time"
)

func TestMainGetsCommand(t *testing.T) {
	t.Parallel()

	callsChan := make(chan string)
	reportingArgsChan := make(chan bool)
	exitArgsChan := make(chan bool)
	mutateResult := make(chan bool)

	// Given a mutation func
	mutate := func() bool {
		callsChan <- "mutate"

		return <-mutateResult
	}

	// Given a reporting func
	report := func(r bool) {
		callsChan <- "report"
		reportingArgsChan <- r
	}

	// Given an exit func
	exit := func(e bool) {
		callsChan <- "exit"
		exitArgsChan <- e
	}

	// When run is called
	go run(mutate, report, exit)

	// Then we run the mutations and return the results
	waitForCall(t, callsChan, "mutate")
	setReturnValue(t, mutateResult, true)

	// Then we report the summary of the run with the mutation results
	waitForCall(t, callsChan, "report")
	waitForArgs(t, reportingArgsChan, true, "report")

	// Then we exit with the result of the mutations
	waitForCall(t, callsChan, "exit")
	waitForArgs(t, exitArgsChan, true, "exit")
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
