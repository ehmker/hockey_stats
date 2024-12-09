package commands

import "github.com/ehmker/hockey_stats/internal/shared"

type CLICommand struct {
	Name string
	Description string
	Callback func(shared.State, []string) error
}

func GetCommands() map[string]CLICommand {
	return map[string]CLICommand{
		"exit": getExitCommand(),
		"update": getScrapeCommand(),
		"fetchtest": getTestFileCommand(),
	}
}

