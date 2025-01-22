-- name: getSourceById :one
SELECT * FROM sources
WHERE id = $1
LIMIT 1;

-- name: listSourcesTemplate :many
SELECT * FROM sources
WHERE id IN ($1, $2)
LIMIT 1;