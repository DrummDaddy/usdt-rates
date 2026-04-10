package grpc_test

import (
	"context"
	"testing"
	"time"

	usdtpb "github.com/DrummDaddy/usdt-rates/gen/gen/usdt/v1"
	"github.com/DrummDaddy/usdt-rates/internal/grpc"
	"github.com/DrummDaddy/usdt-rates/internal/rates/client"
	"github.com/DrummDaddy/usdt-rates/internal/service"
	"github.com/DrummDaddy/usdt-rates/internal/storage/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type fakeClient struct{}

func (f *fakeClient) FetchDepth(ctx context.Context, symbol string) (client.OrderBook, error) {
	return client.OrderBook{
		Asks:      []decimal.Decimal{decimal.NewFromInt(1), decimal.NewFromInt(2), decimal.NewFromInt(3)},
		Bids:      []decimal.Decimal{decimal.NewFromInt(1), decimal.NewFromInt(2), decimal.NewFromInt(3)},
		FetchedAt: time.Now().UTC(),
	}, nil
}

type fakeRepo struct{}

func (r *fakeRepo) SaveRate(ctx context.Context, rec postgres.RateRecord) error {
	return nil
}

func TestHandler_GetRates_InvalidArgument(t *testing.T) {
	svc := service.New(&fakeClient{}, &fakeRepo{})
	h := grpc.NewHandler(svc)

	tests := []struct {
		name string
		n    int32
		m    int32
	}{
		{"n=0", 0, 2},
		{"m=0", 1, 0},
		{"n>m", 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.GetRates(
				context.Background(),
				&usdtpb.GetRatesRequest{N: tt.n, M: tt.m},
			)
			require.Error(t, err)
		})
	}
}
