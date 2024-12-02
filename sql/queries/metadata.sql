-- name: GetLastScrapedDateFromDB :one
SELECT value 
from project_metadata
WHERE key = 'last_scraped_date';

-- name: UpdateLastScrapedDate :exec
UPDATE project_metadata
Set Value = $1
WHERE KEY = 'last_scraped_date';