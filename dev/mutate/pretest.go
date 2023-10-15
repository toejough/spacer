package main

type pretestDeps interface {
	fetchPretestCommand() []string
	runSubprocess([]string) error
	communicateError(error)
}

func pretest(deps pretestDeps) bool {
	command := deps.fetchPretestCommand()

	err := deps.runSubprocess(command)
	if err != nil {
		deps.communicateError(err)
	}

	return err == nil
}
