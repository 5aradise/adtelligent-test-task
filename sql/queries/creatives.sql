-- name: listCreativesByCampaignId :many
SELECT * FROM creatives
WHERE campaign_id = $1;

-- name: listCreativesTemplate :many
SELECT * FROM creatives
WHERE campaign_id IN ($1, $2);