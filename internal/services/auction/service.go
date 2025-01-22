package auction

import (
	"log/slog"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
)

type auctionStorage interface {
}

type service struct {
	stor auctionStorage
	l    *slog.Logger
}

func New(storage auctionStorage, logger *slog.Logger) *service {
	return &service{
		stor: storage,
		l:    logger,
	}
}

func (s *service) GetProfitCreative(sourceId int, maxDuration time.Duration) models.Creative {
	return models.Creative{}
}
