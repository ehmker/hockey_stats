-- name: ResetGameResults :exec
DELETE from game_results;

-- name: ResetPenSummaries :exec
DELETE from penalty_summaries;

-- name: ResetScoringSummaries :exec
DELETE from scoring_summaries;

-- name: ResetSkaterGameStats :exec
DELETE from skater_game_stats;

-- name: ResetShots :exec
DELETE from shots;

-- name: ResetGoalieStats :exec
DELETE FROM goalie_game_stats;