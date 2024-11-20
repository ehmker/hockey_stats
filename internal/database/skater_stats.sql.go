// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: skater_stats.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSkaterGameStats = `-- name: CreateSkaterGameStats :one
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
RETURNING id, gameid, created_at, updated_at, team, player_name, playerid, goals, assists, points, plus_minus, pen_mins, goals_ev, goals_pp, goals_sh, goals_gw, assists_ev, assists_pp, assists_sh, shots, shot_percent, shifts, time_on_ice
`

type CreateSkaterGameStatsParams struct {
	ID         uuid.UUID
	Gameid     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Team       string
	PlayerName string
	Playerid   string
	Goals      int32
	Assists    int32
	Points     int32
	PlusMinus  int32
	PenMins    int32
	GoalsEv    int32
	GoalsPp    int32
	GoalsSh    int32
	GoalsGw    int32
	AssistsEv  int32
	AssistsPp  int32
	AssistsSh  int32
	Shots      int32
	Shifts     int32
	TimeOnIce  int32
}

func (q *Queries) CreateSkaterGameStats(ctx context.Context, arg CreateSkaterGameStatsParams) (SkaterGameStat, error) {
	row := q.db.QueryRowContext(ctx, createSkaterGameStats,
		arg.ID,
		arg.Gameid,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Team,
		arg.PlayerName,
		arg.Playerid,
		arg.Goals,
		arg.Assists,
		arg.Points,
		arg.PlusMinus,
		arg.PenMins,
		arg.GoalsEv,
		arg.GoalsPp,
		arg.GoalsSh,
		arg.GoalsGw,
		arg.AssistsEv,
		arg.AssistsPp,
		arg.AssistsSh,
		arg.Shots,
		arg.Shifts,
		arg.TimeOnIce,
	)
	var i SkaterGameStat
	err := row.Scan(
		&i.ID,
		&i.Gameid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Team,
		&i.PlayerName,
		&i.Playerid,
		&i.Goals,
		&i.Assists,
		&i.Points,
		&i.PlusMinus,
		&i.PenMins,
		&i.GoalsEv,
		&i.GoalsPp,
		&i.GoalsSh,
		&i.GoalsGw,
		&i.AssistsEv,
		&i.AssistsPp,
		&i.AssistsSh,
		&i.Shots,
		&i.ShotPercent,
		&i.Shifts,
		&i.TimeOnIce,
	)
	return i, err
}