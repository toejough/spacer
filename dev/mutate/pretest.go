package main

type pretestDeps interface {
	printStarting(string) func(string)
	fetchPretestCommand() []string
	runSubprocess([]string) bool
}

func pretest(deps pretestDeps) bool {
	done := deps.printStarting("Pretest")
	command := deps.fetchPretestCommand()
	success := deps.runSubprocess(command)

	if success {
		done("Success")
	} else {
		done("Failure")
	}

	return success
}
