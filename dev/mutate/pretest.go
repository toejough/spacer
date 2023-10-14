package main

type pretestDeps interface {
	fetchPretestCommand() []string
	runSubprocess([]string) bool
}

func pretest(deps pretestDeps) bool {
	command := deps.fetchPretestCommand()
	success := deps.runSubprocess(command)

	return success
}
