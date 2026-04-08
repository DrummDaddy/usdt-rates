package service

import (
	"context"
	"time"

	"github.com/DrummDaddy/usdt-rates/internal/rates"
	"github.com/DrummDaddy/usdt-rates/internal/rates/client"
	"github.com/DrummDaddy/usdt-rates/internal/storage/postgres"
	"github.com/shopspring/decimal"
)

const fixedUSDTDepthSymbol = "usdta7a5"

type Service struct {
	client client.GrinexClient
	repo   postgres.Repository
}

func New(grinex client.GrinexClient, repo postgres.Repository) *Service {
	return &Service{client: grinex, repo: repo}
}

type GetRatesInput struct {
	N int
	M int
}
type GetRatesOutput struct {
	FetchedAt time.Time

	AskTopN  decimal.Decimal
	AskAvgNM decimal.Decimal
	BidTopN  decimal.Decimal
	BidAvgNM decimal.Decimal
}

func (s *Service) GetRates(ctx context.Context, in GetRatesInput) (GetRatesOutput, error) {
	ob, err := s.client.FetchDepth(ctx, fixedUSDTDepthSymbol)
	if err != nil {
		return GetRatesOutput{}, err
	}
	askTopN, err := rates.TopN(ob.Asks, in.N)
	if err != nil {
		return GetRatesOutput{}, err
	}
	askAvgNM, err := rates.AvgNM(ob.Asks, in.N, in.M)
	if err != nil {
		return GetRatesOutput{}, err
	}
	bidTopN, err := rates.TopN(ob.Bids, in.N)
	if err != nil {
		return GetRatesOutput{}, err
	}
	bidAvgNM, err := rates.AvgNM(ob.Bids, in.N, in.M)
	if err != nil {
		return GetRatesOutput{}, err
	}

	out := GetRatesOutput{
		FetchedAt: ob.FetchedAt.UTC(),
		AskTopN:   askTopN,
		AskAvgNM:  askAvgNM,
		BidTopN:   bidTopN,
		BidAvgNM:  bidAvgNM,
	}
	if err := s.repo.SaveRate(ctx, postgres.RateRecord{
		FetchedAt: out.FetchedAt,
		AskTopN:   out.AskTopN,
		AskAvgNM:  out.AskAvgNM,
		BidTopN:   out.BidTopN,
		BidAvgNM:  out.BidAvgNM,
	}); err != nil {
		return GetRatesOutput{}, err
	}
	return out, nil
}
