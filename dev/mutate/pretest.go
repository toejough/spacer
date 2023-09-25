package main

type pretestDeps interface {
	printStarting(string) func(string)
	fetchPretestCommand() []string
	runSubprocess([]string)
}

func pretest(deps pretestDeps) bool {
	done := deps.printStarting("Pretest")
	defer done("Success")

	command := deps.fetchPretestCommand()
	deps.runSubprocess(command)

	return true
}
