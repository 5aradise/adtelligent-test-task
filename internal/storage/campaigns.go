package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/ops"
)

func (s *storage) GetCampaignById(ctx context.Context, id int) (models.Campaign, error) {
	const op = "storage.GetCampaignById"

	campaign, ok := s.campaignsCache.Load(id)
	if ok {
		return campaign, nil
	}

	var err error
	campaign, err = s.getCampaignById(ctx, id)
	if err != nil {
		return models.Campaign{}, ops.Wrap(op, err)
	}

	s.campaignsCache.Store(id, campaign)
	return campaign, nil
}

func (s *storage) ListCampaignsByIds(ctx context.Context, ids []int) ([]models.Campaign, error) {
	const op = "storage.ListCampaignsByIds"

	campaigns := make([]models.Campaign, 0, len(ids))
	var idsToGetFromDB []int
	for _, id := range ids {
		campaign, ok := s.campaignsCache.Load(id)
		if ok {
			campaigns = append(campaigns, campaign)
		} else {
			idsToGetFromDB = append(idsToGetFromDB, id)
		}
	}

	fromDBs, err := s.listCampaignsByIds(idsToGetFromDB)
	if err != nil {
		return nil, ops.Wrap(op, err)
	}

	for id, fromDB := range fromDBs {
		campaigns = append(campaigns, fromDB)
		s.campaignsCache.Store(id, fromDB)
	}

	return campaigns, nil
}

const getCampaignById = `
SELECT id, name, start_time, end_time FROM campaigns
WHERE id = $1
LIMIT 1
`

func (s *storage) getCampaignById(ctx context.Context, id int) (models.Campaign, error) {
	const op = "storage.getCampaignById"

	row := s.db.QueryRowContext(ctx, getCampaignById, id)
	var c models.Campaign
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.StartTime,
		&c.EndTime,
	)
	if err != nil {
		return models.Campaign{}, ops.Wrap(op, err)
	}
	return c, nil
}

const listCampaignsTemplate = `
SELECT id, name, start_time, end_time FROM campaigns
WHERE id IN (%s) 
`

func (s *storage) listCampaignsByIds(ids []int) (map[int]models.Campaign, error) {
	const op = "storage.listCampaignsByIds"

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
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(listCampaignsTemplate, strIds), anyIds...)
	if err != nil {
		return nil, ops.Wrap(op, err)
	}
	defer rows.Close()
	campaigns := make(map[int]models.Campaign, len(ids))
	for rows.Next() {
		var campaign models.Campaign
		if err := rows.Scan(
			&campaign.ID,
			&campaign.Name,
			&campaign.StartTime,
			&campaign.EndTime,
		); err != nil {
			return nil, ops.Wrap(op, err)
		}
		campaigns[campaign.ID] = campaign
	}
	if err := rows.Close(); err != nil {
		return nil, ops.Wrap(op, err)
	}
	if err := rows.Err(); err != nil {
		return nil, ops.Wrap(op, err)
	}

	return campaigns, nil
}
