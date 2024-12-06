// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type GameResult struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DatePlayed     time.Time
	Arena          string
	Attendance     int32
	Duration       time.Time
	HomeTeam       string
	HomeTeamScore  int32
	HomeTeamResult string
	HomeTeamRecord string
	AwayTeam       string
	AwayTeamScore  int32
	AwayTeamResult string
	AwayTeamRecord string
}

type GoalieGameStat struct {
	ID           uuid.UUID
	Gameid       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Team         string
	PlayerName   string
	Playerid     string
	Decision     sql.NullString
	GoalsAgainst int32
	ShotsAgainst int32
	Saves        int32
	SavePercent  sql.NullString
	Shutout      bool
	PenMins      int32
	TimeOnIce    int32
}

type PenaltySummary struct {
	ID        uuid.UUID
	Gameid    string
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

type ProjectMetadatum struct {
	Key   string
	Value string
}

type ScoringSummary struct {
	ID              uuid.UUID
	Gameid          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Period          string
	Time            time.Time
	Team            string
	ScoringPlayer   string
	ScoringPlayerID string
	FirstAssist     sql.NullString
	FirstAssistID   sql.NullString
	SecondAssist    sql.NullString
	SecondAssistID  sql.NullString
	EmptyNet        bool
}

type Shot struct {
	ID        uuid.UUID
	Gameid    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Team      string
	Player    string
	XLoc      int32
	YLoc      int32
	Goal      bool
}

type SkaterGameStat struct {
	ID          uuid.UUID
	Gameid      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Team        string
	PlayerName  string
	Playerid    string
	Goals       int32
	Assists     int32
	Points      int32
	PlusMinus   int32
	PenMins     int32
	GoalsEv     int32
	GoalsPp     int32
	GoalsSh     int32
	GoalsGw     int32
	AssistsEv   int32
	AssistsPp   int32
	AssistsSh   int32
	Shots       int32
	ShotPercent sql.NullString
	Shifts      int32
	TimeOnIce   int32
}

type Team struct {
	FullName  string
	ShortName string
}
