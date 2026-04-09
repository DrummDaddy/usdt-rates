package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Asks      []decimal.Decimal
	Bids      []decimal.Decimal
	FetchedAt time.Time
}

type GrinexClient interface {
	FetchDepth(ctx context.Context, symbol string) (OrderBook, error)
}

type grinexClient struct {
	depthURL string
	http     *resty.Client
}

func NewGrinexClient(depthURL string, timeout time.Duration) GrinexClient {
	return &grinexClient{
		depthURL: depthURL,
		http: resty.New().
			SetTimeout(timeout).
			SetRetryCount(1),
	}
}

type depthLevel struct {
	Price string `json:"price"`
}

type depthResponse struct {
	Timestamp int64        `json:"timestamp"`
	Asks      []depthLevel `json:"asks"`
	Bids      []depthLevel `json:"bids"`
}

func (c *grinexClient) FetchDepth(ctx context.Context, symbol string) (OrderBook, error) {
	u, err := url.Parse(c.depthURL)
	if err != nil {
		return OrderBook{}, fmt.Errorf("parse depthURL: %w", err)
	}

	q := u.Query()
	q.Set("symbol", symbol)
	u.RawQuery = q.Encode()

	resp, err := c.http.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get(u.String())
	if err != nil {
		return OrderBook{}, fmt.Errorf("grinex request failed: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return OrderBook{}, fmt.Errorf("grinex http status %d: %s", resp.StatusCode(), resp.String())
	}

	var parsed depthResponse
	if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
		return OrderBook{}, fmt.Errorf("decode grinex json: %w", err)
	}

	toSide := func(in []depthLevel) ([]decimal.Decimal, error) {
		out := make([]decimal.Decimal, 0, len(in))
		for _, lv := range in {
			d, err := decimal.NewFromString(lv.Price)
			if err != nil {
				return nil, fmt.Errorf("parse price %q: %w", lv.Price, err)
			}
			out = append(out, d)
		}
		return out, nil
	}

	asks, err := toSide(parsed.Asks)
	if err != nil {
		return OrderBook{}, fmt.Errorf("parse asks: %w", err)
	}

	bids, err := toSide(parsed.Bids)
	if err != nil {
		return OrderBook{}, fmt.Errorf("parse bids: %w", err)
	}

	return OrderBook{
		Asks:      asks,
		Bids:      bids,
		FetchedAt: time.Unix(parsed.Timestamp, 0).UTC(),
	}, nil
}
