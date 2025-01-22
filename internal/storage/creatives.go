package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

func (s *storage) ListCreativesByCampaignId(ctx context.Context, campaignID int) ([]models.Creative, error) {
	const op = "storage.ListCreativesByCampaignId"

	creatives, ok := s.creativesCache.Load(campaignID)
	if ok {
		return creatives, nil
	}

	var err error
	creatives, err = s.listCreativesByCampaignId(ctx, campaignID)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}

	s.creativesCache.Store(campaignID, creatives)
	return creatives, nil
}

const listCreativesByCampaignId = `
SELECT id, campaign_id, cents_per_view, duration_in_sec, hls_playlist_path FROM creatives
WHERE campaign_id = $1
`

func (s *storage) listCreativesByCampaignId(ctx context.Context, campaignID int) ([]models.Creative, error) {
	const op = "storage.listCreativesByCampaignId"

	rows, err := s.db.QueryContext(ctx, listCreativesByCampaignId, campaignID)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}
	defer rows.Close()
	var creatives []models.Creative
	for rows.Next() {
		var c models.Creative
		if err := rows.Scan(
			&c.ID,
			&c.CampaignID,
			&c.Price,
			&c.Duration,
			&c.HlsPlaylistPath,
		); err != nil {
			return nil, util.OpWrap(op, err)
		}
		creatives = append(creatives, c)
	}
	if err := rows.Close(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	if err := rows.Err(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	return creatives, nil
}

const listCreativesTemplate = `
SELECT id, campaign_id, cents_per_view, duration_in_sec, hls_playlist_path FROM creatives
WHERE campaign_id IN (%s) 
`

func (s *storage) listCreativesByCampaignIds(ids []int) (map[int][]models.Creative, error) {
	const op = "storage.listCreativesByCampaignIds"

	if len(ids) == 0 {
		return nil, nil
	}
	var strIds string
	anyIds := make([]any, len(ids))
	for i, id := range ids {
		strIds += "$" + strconv.Itoa(i+1) + ", "
		anyIds[i] = id
	}
	strIds = strIds[:len(strIds)-2]

	ctx, cancel := context.WithTimeout(context.Background(), s.bigRequestTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(listCreativesTemplate, strIds), anyIds...)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}
	defer rows.Close()
	creatives := make(map[int][]models.Creative, len(ids))
	for rows.Next() {
		var creative models.Creative
		if err := rows.Scan(
			&creative.ID,
			&creative.CampaignID,
			&creative.Price,
			&creative.Duration,
			&creative.HlsPlaylistPath,
		); err != nil {
			return nil, util.OpWrap(op, err)
		}
		creatives[creative.CampaignID] = append(creatives[creative.CampaignID], creative)
	}
	if err := rows.Close(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	if err := rows.Err(); err != nil {
		return nil, util.OpWrap(op, err)
	}

	return creatives, nil
}
