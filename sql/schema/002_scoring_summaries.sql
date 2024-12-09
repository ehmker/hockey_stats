-- +goose Up
CREATE TABLE IF NOT EXISTS scoring_summaries (
    id UUID PRIMARY KEY,
    game_id TEXT NOT NULL REFERENCES game_results(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    period TEXT NOT NULL,
    time TIME NOT NULL, 
    team TEXT NOT NULL,
    player TEXT NOT NULL,
    player_id TEXT NOT NULL,
    first_assist TEXT,
    first_assist_id TEXT,
    second_assist TEXT, 
    second_assist_id TEXT, 
    empty_net BOOLEAN NOT NULL,

    CONSTRAINT unique_scoring_summary UNIQUE (player_id, game_id, time)
);

-- +goose Down
DROP TABLE scoring_summaries;