package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type RateRecord struct {
	FetchedAt time.Time
	AskTopN   decimal.Decimal
	AskAvgNM  decimal.Decimal
	BidTopN   decimal.Decimal
	BidAvgNM  decimal.Decimal
}

type Repository interface {
	SaveRate(ctx context.Context, r RateRecord) error
}

type repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return &repo{pool: pool}
}

func (r *repo) SaveRate(ctx context.Context, rec RateRecord) error {
	if rec.AskTopN.IsNegative() {
		return fmt.Errorf("ask_top_n cannot be negative")
	}
	_, err := r.pool.Exec(ctx, `
INSERT INTO rates (fetched_ad, ask_avg_n_m, bid_top_n, big_avg_n_m)
VALUES ($1, $2, $3, $4)
`, rec.FetchedAt, rec.AskTopN, rec.AskAvgNM, rec.BidTopN, rec.BidAvgNM)
	if err != nil {
		return fmt.Errorf(" failed to save rate: %w", err)
	}
	return nil
}
