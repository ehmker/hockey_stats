// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: penalty_summary.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPenaltySummary = `-- name: CreatePenaltySummary :one
INSERT INTO
    penalty_summaries (
        id,
        game_id,
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
RETURNING id, game_id, created_at, updated_at, period, time, team, player, player_id, penalty, pim
`

type CreatePenaltySummaryParams struct {
	ID        uuid.UUID
	GameID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Period    string
	Time      time.Time
	Team      string
	Player    string
	PlayerID  string
	Penalty   string
	Pim       int32
}

func (q *Queries) CreatePenaltySummary(ctx context.Context, arg CreatePenaltySummaryParams) (PenaltySummary, error) {
	row := q.db.QueryRowContext(ctx, createPenaltySummary,
		arg.ID,
		arg.GameID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Period,
		arg.Time,
		arg.Team,
		arg.Player,
		arg.PlayerID,
		arg.Penalty,
		arg.Pim,
	)
	var i PenaltySummary
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Period,
		&i.Time,
		&i.Team,
		&i.Player,
		&i.PlayerID,
		&i.Penalty,
		&i.Pim,
	)
	return i, err
}
