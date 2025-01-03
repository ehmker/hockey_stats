-- name: CreateGoalieStats :one

INSERT INTO GOALIE_GAME_STATS (
    ID,
    game_id, 
    CREATED_AT,
    UPDATED_AT,
    TEAM,
    PLAYER_NAME,
    player_id,
    DECISION,
    GOALS_AGAINST,
    SHOTS_AGAINST,
    SAVES, 
    SHUTOUT,
    PEN_MINS,
    TIME_ON_ICE,
    SEASON
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;
