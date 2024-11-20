-- name: CreatePenaltySummary :one
INSERT INTO
    penalty_summaries (
        id,
        gameid,
        created_at,
        updated_at,
        period,
        time, 
        team,
        player,
        player_id,
        penalty,
        pim
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;