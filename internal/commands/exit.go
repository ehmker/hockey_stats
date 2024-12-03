package commands

import (
	"fmt"
	"os"

	"github.com/ehmker/hockey_stats/internal/shared"
)

func getExitCommand() CLICommand{
	return CLICommand{
		Name: "exit",
		Description: "exits the program",
		Callback: commandExit,
	}
}

func commandExit(_ shared.State, _ string) error{
	fmt.Println("Exiting.")
	os.Exit(0)
	return nil
}