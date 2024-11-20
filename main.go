package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ehmker/hockey_stats/internal/config"
	"github.com/ehmker/hockey_stats/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db *database.Queries
}
func CreateState() (state, error) {
	c := config.Read()
	db, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		return state{}, fmt.Errorf("problem opening database: %v", err)
	}
	dbQueries := database.New(db)
	return state{
		cfg: &c,
		db: dbQueries,
	}, nil
}

func (s state) resetDB() error {
	err := s.db.ResetSkaterGameStats(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetPenSummaries(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetScoringSummaries(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetGameResults(context.Background())
	if err != nil {
		return err
	}
	return nil
}	

func main() {
	s, err := CreateState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.resetDB()
	if err != nil {
		log.Fatal(err)
	}
	gameResultParams, err := ScrapeGameResults()
	if err != nil{
		log.Fatalf("error scraping game results %v", err)
	}

	_, err = s.db.CreateGameResult(context.Background(), gameResultParams)
	if err != nil {
		log.Fatalln(err)
	}

	AddPenaltySummary(s)
	AddScoringSummaryToDB(s)
	AddPlayerStats(s)
}

