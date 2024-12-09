package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/ehmker/hockey_stats/internal/web_scraping"
)

func getScrapeCommand() CLICommand {
	return CLICommand{
		Name: "scrape",
		Description: "Fetches and scrapes game data from hockey-reference.com",
		Callback: ScrapeData,
	}
}

func ScrapeData(s shared.State, _ []string) error {
	startDate := s.GetLastScrapedDate().AddDate(0, 0, 1) //Start date is day after last scraped
	today := time.Now()
	// End date is set to yesterday as site does not update with scores same day
	yesterday := time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, today.Location())  // setting time to midnight to normalize comparisons
	for gameday := startDate; gameday.Before(yesterday); gameday = gameday.AddDate(0, 0, 1){
		addGamesOfDay(s, gameday.Format("2006/01/02"))
	}
	s.SetLastScrapedDate()
	return nil
}

// game_date structure = "YYYY/MM/DD"
func addGamesOfDay(s shared.State, game_date string) {
	url := "https://www.hockey-reference.com/boxscores/"
	if game_date != "" {
		url += game_date
		fmt.Printf("getting games played on: %v\n", game_date)
	}

	gamesList := web_scraping.ScrapeGameLinks(url)

	for _, game := range gamesList {
		gameResult, err := web_scraping.AddGameToDB(s, game)
		if err != nil {
			log.Printf("error adding game to database:\n\t%v\n", err)
			return
		}
		printGame(gameResult)
		time.Sleep(20 * time.Second) // 
		
		
	}


}


func printGame(game database.CreateGameResultParams) {
	fmt.Printf("%s | (%s) %s: %v at (%s) %s: %v\n", game.ID, game.AwayTeamResult, game.AwayTeam, game.AwayTeamScore, game.HomeTeamResult, game.HomeTeam, game.HomeTeamScore)
}