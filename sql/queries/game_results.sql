-- name: CreateGameResult :one
INSERT INTO
    game_results (
        id,
        created_at, 
        updated_at, 
        season,
        date_played, 
        arena,
        attendance,
        duration,
        home_team, 
        home_team_score, 
        home_team_result, 
        home_team_record,
        away_team, 
        away_team_score, 
        away_team_result, 
        away_team_record
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: GetDateOfLastResult :one
SELECT date_played
From game_results
ORDER BY date_played DESC
LIMIT 1;