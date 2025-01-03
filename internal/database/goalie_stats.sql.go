// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: goalie_stats.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createGoalieStats = `-- name: CreateGoalieStats :one

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
RETURNING id, game_id, created_at, updated_at, team, season, player_name, player_id, decision, goals_against, shots_against, saves, save_percent, shutout, pen_mins, time_on_ice
`

type CreateGoalieStatsParams struct {
	ID           uuid.UUID
	GameID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Team         string
	PlayerName   string
	PlayerID     string
	Decision     sql.NullString
	GoalsAgainst int32
	ShotsAgainst int32
	Saves        int32
	Shutout      int32
	PenMins      int32
	TimeOnIce    int32
	Season       int32
}

func (q *Queries) CreateGoalieStats(ctx context.Context, arg CreateGoalieStatsParams) (GoalieGameStat, error) {
	row := q.db.QueryRowContext(ctx, createGoalieStats,
		arg.ID,
		arg.GameID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Team,
		arg.PlayerName,
		arg.PlayerID,
		arg.Decision,
		arg.GoalsAgainst,
		arg.ShotsAgainst,
		arg.Saves,
		arg.Shutout,
		arg.PenMins,
		arg.TimeOnIce,
		arg.Season,
	)
	var i GoalieGameStat
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Team,
		&i.Season,
		&i.PlayerName,
		&i.PlayerID,
		&i.Decision,
		&i.GoalsAgainst,
		&i.ShotsAgainst,
		&i.Saves,
		&i.SavePercent,
		&i.Shutout,
		&i.PenMins,
		&i.TimeOnIce,
	)
	return i, err
}
