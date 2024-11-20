-- name: CreateSkaterGameStats :one
INSERT INTO skater_game_stats (
    ID,
    GAMEID, 
    CREATED_AT,
    UPDATED_AT,
    TEAM, 
    PLAYER_NAME, 
    PLAYERID, 
    GOALS, 
    ASSISTS,
    POINTS, 
    PLUS_MINUS,
    PEN_MINS,
    GOALS_EV,
    GOALS_PP,
    GOALS_SH,
    GOALS_GW,
    ASSISTS_EV,
    ASSISTS_PP,
    ASSISTS_SH,
    SHOTS,
    SHIFTS,
    TIME_ON_ICE
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
RETURNING *;