package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/cache"
)

type storage struct {
	db                DBTX
	bigRequestTimeout time.Duration
	sourcesCache      cache.Cache[int, models.Source]
	campaignsCache    cache.Cache[int, models.Campaign]
	creativesCache    cache.Cache[int, []models.Creative]
}

type DBTX interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func New(db DBTX, logger *slog.Logger, bigRequestTimeout, updateCacheDelay time.Duration) *storage {
	s := &storage{
		db:                db,
		bigRequestTimeout: bigRequestTimeout,
	}
	s.sourcesCache = cache.New(logger.With(slog.String("cache", "sources")), s.listSourcesByIds, updateCacheDelay)
	s.campaignsCache = cache.New(logger.With(slog.String("cache", "campaigns")), s.listCampaignsByIds, updateCacheDelay)
	s.creativesCache = cache.New(logger.With(slog.String("cache", "creatives")), s.listCreativesByCampaignIds, updateCacheDelay)
	return s
}
