package commands

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ehmker/hockey_stats/internal/shared"
)

func getTestFileCommand() CLICommand{
	return CLICommand{
		Name: "fetch-test-file",
		Description: "saves raw html of request to test_files folder for debugging purposes.\n usage: testfetch <gameid> <file name>",
		Callback: FetchTestGameFile,
	}
}

func FetchTestGameFile(s shared.State, inputParams []string) error {
	gameid := strings.ToUpper(inputParams[0])
	url := "https://www.hockey-reference.com/boxscores/" + gameid + ".html"
	filepath := "test_files/"+ gameid + "_" +inputParams[1]+".htm"
	log.Printf("requesting: %v", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed with status code: %d\n", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Failed to create file: %v\n", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save response to file: %v\n", err)
	}

	log.Printf("HTML response saved to: %v\n", filepath)
	return nil
}