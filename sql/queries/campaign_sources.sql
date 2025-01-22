-- name: listCampaignIdBySourceId :many
SELECT campaign_id
FROM campaign_sources
WHERE source_id = $1;