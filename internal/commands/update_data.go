package commands

import (
	"fmt"
	"time"

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

func ScrapeData(s shared.State, _ string) error {
	startDate := s.GetLastScrapedDate().AddDate(0, 0, 1) //Start date is day after last scraped
	endDate := time.Now().AddDate(0, 0, -1) // End date is set to yesterday as site does not update with scores same day
	
	for gameday := startDate; gameday.Before(endDate); gameday = gameday.AddDate(0, 0, 1){
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
		web_scraping.AddGameToDB(s, game)
		fmt.Printf("url: %v | ID: %v\n", game.Url, game.Gameid)
		time.Sleep(20 * time.Second)
	}


}
