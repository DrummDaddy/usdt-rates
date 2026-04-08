package client

import (
	"context"
	"testing"
	"time"

	//"github.com/DrummDaddy/usdt-rates/internal/rates/client"
	_ "github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFetchDepth_ParsesTimestampAndSides(t *testing.T) {
	const depthURL = "https://grinex.io/api/v1/spot/depth"

	jsonBody := `{

	  "timestamp": 1775627101,
	  "asks": [
		{"price":"80.73","volume":"19047.52","amount":"1537706.29"}
	  ],
	  "bids": [
		{"price":"79.73","volume":"19047.52","amount":"1537706.29"}
	  ]
	}`
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		"https://grinex.io/api/v1/spot/depth?symbol=usdta7a5",
		httpmock.NewStringResponder(200, jsonBody),
	)
	c := NewGrinexClient(depthURL, 2*time.Second).(interface {
		FetchDepth(ctx context.Context, symbol string) (OrderBook, error)
	})
	ob, err := c.FetchDepth(context.Background(), "usdta7a5")
	require.NoError(t, err)

	require.Equal(t, time.Unix(1775627101, 0).UTC(), ob.FetchedAt)
	require.Equal(t, []decimal.Decimal{decimal.RequireFromString("80.73")}, ob.Asks)
	require.Equal(t, []decimal.Decimal{decimal.RequireFromString("79.73")}, ob.Bids)
}
