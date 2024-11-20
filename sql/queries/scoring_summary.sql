-- name: CreateScoringSummary :one
INSERT INTO
    scoring_summaries (
        id,
        gameid,
        created_at,
        updated_at,
        period,
        time, 
        team,
        scoring_player,
        scoring_player_id,
        first_assist,
        first_assist_id,
        second_assist, 
        second_assist_id,
        empty_net
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;