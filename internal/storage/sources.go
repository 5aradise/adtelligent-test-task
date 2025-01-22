package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

func (s *storage) GetSourceById(ctx context.Context, id int) (models.Source, error) {
	const op = "storage.GetSourceById"

	source, ok := s.sourcesCache.Load(id)
	if ok {
		return source, nil
	}

	var err error
	source, err = s.getSourceById(ctx, id)
	if err != nil {
		return models.Source{}, util.OpWrap(op, err)
	}

	s.sourcesCache.Store(id, source)
	return source, nil
}

const getSourceById = `
SELECT id, name, is_active FROM sources
WHERE id = $1
LIMIT 1
`

func (s *storage) getSourceById(ctx context.Context, id int) (models.Source, error) {
	const op = "storage.getSourceById"

	var source models.Source

	row := s.db.QueryRowContext(ctx, getSourceById, id)
	err := row.Scan(&source.ID, &source.Name, &source.IsActive)
	if err != nil {
		return models.Source{}, util.OpWrap(op, err)
	}

	source.CampaignIds, err = s.listCampaignIdsBySourceId(ctx, id)
	if err != nil {
		return models.Source{}, util.OpWrap(op, err)
	}
	return source, nil
}

const listSourcesTemplate = `
SELECT id, name, is_active FROM sources
WHERE id IN (%s) 
`

func (s *storage) listSourcesByIds(ids []int) (map[int]models.Source, error) {
	const op = "storage.listSourcesByIds"

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
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(listSourcesTemplate, strIds), anyIds...)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}
	defer rows.Close()
	sources := make(map[int]models.Source, len(ids))
	for rows.Next() {
		var source models.Source
		if err := rows.Scan(&source.ID, &source.Name, &source.IsActive); err != nil {
			return nil, util.OpWrap(op, err)
		}
		sources[source.ID] = source
	}
	if err := rows.Close(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	if err := rows.Err(); err != nil {
		return nil, util.OpWrap(op, err)
	}

	for id, source := range sources {
		source.CampaignIds, err = s.listCampaignIdsBySourceId(context.Background(), id)
		if err != nil {
			return nil, util.OpWrap(op, err)
		}
		sources[id] = source
	}

	return sources, nil
}

const listCampaignIdsBySourceId = `
SELECT campaign_id
FROM campaign_sources
WHERE source_id = $1
`

func (s *storage) listCampaignIdsBySourceId(ctx context.Context, id int) ([]int, error) {
	const op = "storage.listCampaignIdsBySourceId"

	var campaignIds []int

	rows, err := s.db.QueryContext(ctx, listCampaignIdsBySourceId, id)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var campaignId int
		if err := rows.Scan(&campaignId); err != nil {
			return nil, util.OpWrap(op, err)
		}
		campaignIds = append(campaignIds, campaignId)
	}

	if err := rows.Close(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	if err := rows.Err(); err != nil {
		return nil, util.OpWrap(op, err)
	}
	return campaignIds, nil
}
