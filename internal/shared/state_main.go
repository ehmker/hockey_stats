package shared

import (
	"database/sql"
	"fmt"

	"github.com/ehmker/hockey_stats/internal/config"
	"github.com/ehmker/hockey_stats/internal/database"
)

type State struct {
	Cfg *config.Config
	DB  *database.Queries
}



func CreateState() (State, error) {
	c := config.Read()
	db, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		return State{}, fmt.Errorf("problem opening database: %v", err)
	}
	dbQueries := database.New(db)
	return State{
		Cfg: &c,
		DB: dbQueries,
	}, nil
}