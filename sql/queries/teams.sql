-- name: GetTeamShortName :one
SELECT short_name 
FROM   teams
where  full_name = $1;