package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ehmker/hockey_stats/internal/commands"
	"github.com/ehmker/hockey_stats/internal/shared"
)

func startREPL() {
	reader := bufio.NewScanner(os.Stdin)
	commands := commands.GetCommands()
	s, err := shared.CreateState()
	if err != nil{
		log.Printf("error creating state: %v\n", err)
		os.Exit(1)
	}
	last_scraped, err := s.DB.GetLastScrapedDateFromDB(context.Background())
	if err != nil {
		log.Printf("error getting last scraped date: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Last scraped on: %v\n", last_scraped)

	for {
		fmt.Print("Stats > ")
		reader.Scan()
		inputCMD, inputLocation := processInput(reader.Text())

		if cmd, ok := commands[inputCMD]; ok {
			cmd.Callback(s, inputLocation)
		} else {
			fmt.Printf("'%v' command not found\n", inputCMD)
		}
	}

}

func processInput(user_input string) (command, location string) {
	user_input = strings.ToLower(user_input)
	split_input := strings.Split(user_input, " ")

	if len(split_input) == 2{
		return split_input[0], split_input[1]
	}
	return split_input[0], ""


}