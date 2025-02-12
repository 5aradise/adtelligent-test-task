package auction

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/slices"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

var (
	ErrInactiveSource  = errors.New("inactive source")
	ErrCreativeUnfound = errors.New("creative unfound")
)

type auctionStorage interface {
	GetSourceById(ctx context.Context, id int) (models.Source, error)
	ListCampaignsByIds(ctx context.Context, ids []int) ([]models.Campaign, error)
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

func (s *service) GetProfitCreative(sourceId int, maxDuration time.Duration, relativeTime time.Time) (models.Creative, error) {
	const op = "service.GetProfitCreative"
	l := s.l.With(
		slog.String("op", op),
		slog.Int("source_id", sourceId),
	)

	source, err := s.stor.GetSourceById(context.Background(), sourceId)
	if err != nil {
		l.Warn("cannot get source by id", util.SlErr(err))
		return models.Creative{}, util.OpWrap(op, err)
	}
	if err := isValidSource(source); err != nil {
		return models.Creative{}, util.OpWrap(op, err)
	}

	campaigns, err := s.stor.ListCampaignsByIds(context.Background(), source.CampaignIds)
	if err != nil {
		l.Warn("cannot list campaigns by ids", slog.Any("campaigns_id", source.CampaignIds), util.SlErr(err))
		return models.Creative{}, util.OpWrap(op, err)
	}
	bestCreatives := make([]models.Creative, len(campaigns))
	wg := &sync.WaitGroup{}
	for i, campaign := range campaigns {
		if !campaign.IsActive(relativeTime) {
			continue
		}

		wg.Add(1)
		go func() {
			l = l.With(slog.Int("campaign_id", campaign.ID))
			defer wg.Done()

			creatives, err := s.stor.ListCreativesByCampaignId(context.Background(), campaign.ID)
			if err != nil {
				l.Warn("cannot list creatives by campaign id", util.SlErr(err))
				return
			}
			bestCreative, isFound := chooseBestCreative(creatives, maxDuration)
			if isFound {
				bestCreatives[i] = bestCreative
			}
		}()
	}
	wg.Wait()
	bestCreatives = filterEmptyCreatives(bestCreatives)

	bestCreative, isFound := chooseBestCreative(bestCreatives, maxDuration)
	if !isFound {
		return models.Creative{}, util.OpWrap(op, ErrCreativeUnfound)
	}

	return bestCreative, nil
}

func isValidSource(s models.Source) error {
	if !s.IsActive {
		return ErrInactiveSource
	}
	return nil
}

var minPriceCreative = models.Creative{}

func chooseBestCreative(cs []models.Creative, maxDuration time.Duration) (bestCreative models.Creative, isFound bool) {
	if len(cs) < 1 {
		return models.Creative{}, false
	}

	bestCreative = minPriceCreative
	for _, c := range cs {
		creativeDuration := time.Duration(c.DurationInMs) * time.Millisecond
		if creativeDuration > maxDuration {
			continue
		}
		if c.Price > bestCreative.Price {
			bestCreative = c
		}
	}
	if bestCreative.ID == minPriceCreative.ID {
		return models.Creative{}, false
	}
	return bestCreative, true
}

func filterEmptyCreatives(cs []models.Creative) []models.Creative {
	return slices.FilterFunc(cs, func(c models.Creative) bool {
		return c.ID != 0
	})
}
