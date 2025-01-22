-- name: getCampaignById :one
SELECT * FROM campaigns
WHERE id = $1
LIMIT 1;

-- name: listCampaignsTemplate :one
SELECT * FROM campaigns
WHERE id IN ($1, $2);