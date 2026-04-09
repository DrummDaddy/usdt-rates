package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFetchDepth_ParsesTimestampAndSides(t *testing.T) {
	const depthURL = "https://grinex.io/api/v1/spot/depth"
	const symbol = "usdta7a5"

	const ts int64 = 1775717123

	jsonBody := fmt.Sprintf(`{
	  "timestamp": %d,
	  "asks": [
		{"price":"80.73","volume":"19047.52","amount":"1537706.29"}
	  ],
	  "bids": [
		{"price":"79.73","volume":"19047.52","amount":"1537706.29"}
	  ]
	}`, ts)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("%s?symbol=%s", depthURL, symbol),
		httpmock.NewStringResponder(200, jsonBody),
	)

	c := NewGrinexClient(depthURL, 2*time.Second)
	ob, err := c.FetchDepth(context.Background(), symbol)
	require.NoError(t, err)

	require.Equal(t, time.Unix(ts, 0).UTC(), ob.FetchedAt)
	require.Equal(t, []decimal.Decimal{decimal.RequireFromString("80.73")}, ob.Asks)
	require.Equal(t, []decimal.Decimal{decimal.RequireFromString("79.73")}, ob.Bids)
}
