-- +goose Up
CREATE TABLE IF NOT EXISTS scoring_summaries (
    id UUID PRIMARY KEY,
    gameid TEXT NOT NULL REFERENCES game_results(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    period TEXT NOT NULL,
    time TIME NOT NULL, 
    team TEXT NOT NULL,
    scoring_player TEXT NOT NULL,
    scoring_player_id TEXT NOT NULL,
    first_assist TEXT,
    first_assist_id TEXT,
    second_assist TEXT, 
    second_assist_id TEXT, 
    empty_net BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE scoring_summaries;