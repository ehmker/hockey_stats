-- name: CreateShot :one
INSERT INTO shots (
    id, 
    game_id, 
    created_at,
    updated_at,
    team, 
    player, 
    x_loc,
    y_loc,
    goal
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;