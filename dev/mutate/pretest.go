package main

type pretestDeps interface {
	printStarting(string) func(string)
	fetchPretestCommand() []string
	runSubprocess([]string) bool
}

func pretest(deps pretestDeps) bool {
	var success bool

	done := deps.printStarting("Pretest")

	defer func() {
		if success {
			done("Success")
		} else {
			done("Failure")
		}
	}()

	command := deps.fetchPretestCommand()
	success = deps.runSubprocess(command)

	return success
}
