package main

type pretestDeps interface {
	fetchPretestCommand() []string
	runSubprocess([]string) error
}

func pretest(deps pretestDeps) bool {
	command := deps.fetchPretestCommand()
	err := deps.runSubprocess(command)
	// TODO: some sort of notification for the error

	return err == nil
}
