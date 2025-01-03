// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: metadata.sql

package database

import (
	"context"
)

const getLastScrapedDateFromDB = `-- name: GetLastScrapedDateFromDB :one
SELECT value 
from project_metadata
WHERE key = 'last_scraped_date'
`

func (q *Queries) GetLastScrapedDateFromDB(ctx context.Context) (string, error) {
	row := q.db.QueryRowContext(ctx, getLastScrapedDateFromDB)
	var value string
	err := row.Scan(&value)
	return value, err
}

const updateLastScrapedDate = `-- name: UpdateLastScrapedDate :exec
UPDATE project_metadata
Set Value = $1
WHERE KEY = 'last_scraped_date'
`

func (q *Queries) UpdateLastScrapedDate(ctx context.Context, value string) error {
	_, err := q.db.ExecContext(ctx, updateLastScrapedDate, value)
	return err
}
