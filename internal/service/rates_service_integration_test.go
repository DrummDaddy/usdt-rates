package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/DrummDaddy/usdt-rates/internal/rates/client"
	"github.com/DrummDaddy/usdt-rates/internal/service"
	"github.com/DrummDaddy/usdt-rates/internal/storage/postgres"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type fakeClient struct {
	ob client.OrderBook
}

func (f *fakeClient) FetchDepth(ctx context.Context, symbol string) (client.OrderBook, error) {
	return f.ob, nil
}

type spyRepo struct {
	called bool
	rec    postgres.RateRecord
}

func (r *spyRepo) SaveRate(ctx context.Context, rec postgres.RateRecord) error {
	r.called = true
	r.rec = rec
	return nil
}

func TestService_GetRates_HappyPath_SavesRate(t *testing.T) {
	fixedAt := time.Date(2026, 4, 9, 7, 0, 0, 0, time.UTC)

	fake := &fakeClient{
		ob: client.OrderBook{
			Asks:      []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20), decimal.NewFromInt(30), decimal.NewFromInt(40)},
			Bids:      []decimal.Decimal{decimal.NewFromInt(1), decimal.NewFromInt(2), decimal.NewFromInt(3), decimal.NewFromInt(4)},
			FetchedAt: fixedAt,
		},
	}
	repo := &spyRepo{}

	svc := service.New(fake, repo)

	out, err := svc.GetRates(context.Background(), service.GetRatesInput{N: 2, M: 4})
	require.NoError(t, err)

	require.True(t, repo.called)
	require.True(t, repo.rec.FetchedAt.Equal(fixedAt))
	require.True(t, out.FetchedAt.Equal(fixedAt))

	require.True(t, out.AskTopN.Equal(decimal.NewFromInt(20)))
	require.True(t, out.AskAvgNM.Equal(decimal.NewFromInt(30)))

	// bids: [1,2,3,4], n=2,m=4 => topN=2, avg=(2+3+4)/3=3
	require.True(t, out.BidTopN.Equal(decimal.NewFromInt(2)))
	require.True(t, out.BidAvgNM.Equal(decimal.NewFromInt(3)))
}
