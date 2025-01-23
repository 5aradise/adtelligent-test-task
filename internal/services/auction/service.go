package auction

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

var (
	ErrInactiveSource  = errors.New("inactive source")
	ErrCreativeUnfound = errors.New("creative unfound")
)

type auctionStorage interface {
	GetSourceById(ctx context.Context, id int) (models.Source, error)
	GetCampaignById(ctx context.Context, id int) (models.Campaign, error)
	ListCreativesByCampaignId(ctx context.Context, id int) ([]models.Creative, error)
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

var minPriceCreative = models.Creative{}

func (s *service) GetProfitCreative(sourceId int, maxDuration time.Duration) (models.Creative, error) {
	const op = "service.GetProfitCreative"
	l := s.l.With(
		slog.String("op", op),
		slog.Int("source_id", sourceId),
	)

	source, err := s.stor.GetSourceById(context.Background(), sourceId)
	if err != nil {
		return models.Creative{}, util.OpWrap(op, err)
	}

	if !source.IsActive {
		return models.Creative{}, util.OpWrap(op, ErrInactiveSource)
	}

	now := time.Now()
	bestCreative := minPriceCreative
	for _, campaignId := range source.CampaignIds {
		campaign, err := s.stor.GetCampaignById(context.Background(), campaignId)
		l = l.With(slog.Int("campaign_id", campaignId))
		if err != nil {
			l.Warn("cannot get campaign by id", util.SlErr(err))
			continue
		}

		if !(campaign.StartTime.Before(now) && campaign.EndTime.After(now)) {
			l.Debug("campaign doesn't fit into time frame", slog.Time("now", now), slog.Time("start", campaign.StartTime), slog.Time("end", campaign.EndTime))
			continue
		}

		creatives, err := s.stor.ListCreativesByCampaignId(context.Background(), campaign.ID)
		if err != nil {
			l.Warn("cannot list creatives by campaign id", util.SlErr(err))
			continue
		}

		for _, creative := range creatives {
			l = l.With(slog.Int("creative_id", creative.ID))
			creativeDuration := time.Duration(creative.DurationInMs) * time.Millisecond
			if creativeDuration > maxDuration {
				l.Debug("creative takes longer than max duration", slog.Duration("creative_duration", creativeDuration), slog.Duration("max_duration", maxDuration))
				continue
			}

			if creative.Price > bestCreative.Price {
				bestCreative = creative
			}
		}
	}

	if bestCreative.ID == 0 {
		return models.Creative{}, util.OpWrap(op, ErrCreativeUnfound)
	}
	return bestCreative, nil
}
