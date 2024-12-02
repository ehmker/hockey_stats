package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ehmker/hockey_stats/internal/config"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/ehmker/hockey_stats/internal/web_scraping"
	_ "github.com/lib/pq"
)

func CreateState() (shared.State, error) {
	c := config.Read()
	db, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		return shared.State{}, fmt.Errorf("problem opening database: %v", err)
	}
	dbQueries := database.New(db)
	return shared.State{
		Cfg: &c,
		DB: dbQueries,
	}, nil
}



func main() {
	s, err := CreateState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// err = s.ResetDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Open the local HTML file
	// file, err := os.Open("example_pages/example_3.htm")
	// if err != nil {
	// 	log.Fatalf("Error opening file: %v", err)
	// }
	// defer file.Close()

	// // Load the HTML document
	// doc, err := goquery.NewDocumentFromReader(file)
	// if err != nil {
	// 	log.Fatalf("Error parsing HTML: %v", err)
	// }

	// // Step 1: Select the last row in tbody
	// lastRow := doc.Find(".game_summary.nohover.current tbody tr").Last()

	// // Step 2: Select the last td in that row
	// lastCell := lastRow.Find("td").Last().Text()

	// fmt.Println(strings.TrimSpace(lastCell))
	startDate := GetLastScrapedDate(s).AddDate(0, 0, 1) //Start date is day after last scraped
	endDate := time.Now().AddDate(0, 0, -1) // End date is set to yesterday as site does not update with scores same day
	
	for gameday := startDate; gameday.Before(endDate); gameday = gameday.AddDate(0, 0, 1){
		AddGamesOfDay(s, gameday.Format("2006/01/02"))
	}
	SetLastScrapedDate(s)
}

// game_date structure = "YYYY/MM/DD"
func AddGamesOfDay(s shared.State, game_date string) {
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
		// break
	}


}

func GetLastScrapedDate(s shared.State) time.Time {
	lastScraped, err := s.DB.GetLastScrapedDateFromDB(context.Background())
	if err != nil {
		log.Printf("error getting last scraped date from DB: %v\n", err)
		return time.Now()
	}
	
	lastScraped_time, err := time.Parse("2006/01/02", lastScraped)
	if err != nil {
		log.Printf("error parsing last scraped date as time: %v\n", err)
		return time.Now()
	}

	return lastScraped_time
}

func SetLastScrapedDate(s shared.State) {
	lastScrapedGameDate, err := s.DB.GetDateOfLastResult(context.Background())
	if err != nil {
		log.Printf("error getting last game date: %v\n", err)
		lastScrapedGameDate = time.Now()
	}

	err = s.DB.UpdateLastScrapedDate(context.Background(), lastScrapedGameDate.Format("2006/01/02"))
	if err != nil{
		log.Printf("error updating last scraped date in DB: %v\n", err)
	}

}