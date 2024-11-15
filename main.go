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

func main() {
	s, err := CreateState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	gameResultParams, err := ScrapeGameResults()
	if err != nil{
		log.Fatalf("error scraping game results %v", err)
	}

	_, err = s.db.CreateGameResult(context.Background(), gameResultParams)
	if err != nil {
		log.Fatalln(err)
	}
	
}


// func ScrapeBoxScore () {

// }