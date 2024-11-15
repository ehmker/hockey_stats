-- +goose Up
CREATE TABLE IF NOT EXISTS game_results (
    id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    date_played TIMESTAMP NOT NULL, 
    arena TEXT NOT NULL,
    attendance INTEGER NOT NULL,
    duration TIME NOT NULL,
    home_team TEXT NOT NULL, 
    home_team_score integer NOT NULL, 
    home_team_result TEXT NOT NULL, 
    home_team_record TEXT NOT NULL,
    away_team TEXT NOT NULL, 
    away_team_score integer NOT NULL, 
    away_team_result TEXT NOT NULL, 
    away_team_record TEXT NOT NULL
);

-- +goose Down
DROP TABLE game_results;