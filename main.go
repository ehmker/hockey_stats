package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

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
	err = s.ResetDB()
	if err != nil {
		log.Fatal(err)
	}
	gameResultParams, err := web_scraping.ScrapeGameResults()
	if err != nil{
		log.Fatalf("error scraping game results %v", err)
	}

	_, err = s.DB.CreateGameResult(context.Background(), gameResultParams)
	if err != nil {
		log.Fatalln(err)
	}

	web_scraping.AddPenaltySummary(s)
	web_scraping.AddScoringSummaryToDB(s)
	web_scraping.AddPlayerStats(s)
	web_scraping.AddShotLocationsToDB(s)
	web_scraping.AddGoalieStats(s)
}

