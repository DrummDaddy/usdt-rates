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
	if rec.AskAvgNM.IsNegative() {
		return fmt.Errorf("ask_avg_n_m cannot be negative")
	}
	if rec.BidTopN.IsNegative() {
		return fmt.Errorf("bid_top_n cannot be negative")
	}
	if rec.BidAvgNM.IsNegative() {
		return fmt.Errorf("bid_avg_n_m cannot be negative")
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO rates (fetched_at, ask_top_n, ask_avg_n_m, bid_top_n, bid_avg_n_m)
		VALUES ($1, $2, $3, $4, $5)
	`, rec.FetchedAt, rec.AskTopN, rec.AskAvgNM, rec.BidTopN, rec.BidAvgNM)

	if err != nil {
		return fmt.Errorf("failed to save rate: %w", err)
	}

	return nil
}
