-- +goose Up
CREATE TABLE IF NOT EXISTS penalty_summaries (
    id UUID PRIMARY KEY,
    gameid TEXT NOT NULL REFERENCES game_results(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    period TEXT NOT NULL,
    time TIME NOT NULL, 
    team TEXT NOT NULL,
    player TEXT NOT NULL,
    player_id TEXT NOT NULL,
    penalty TEXT NOT NULL,
    pim INTEGER NOT NULL
);

-- +goose Down
DROP TABLE penalty_summaries;