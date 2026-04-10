package postgres_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrummDaddy/usdt-rates/internal/storage/postgres"
	"github.com/shopspring/decimal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func findMigrationsPath(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)

	// Walk up until we find internal/storage/migrations (or migrtations if you didn't rename)
	for i := 0; i < 10; i++ {
		candidate := filepath.Join(wd, "internal", "storage", "migrations")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		// go up one dir
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	t.Fatal("cannot find internal/storage/migrations")
	return ""
}

func TestRepository_SaveRate_WritesToPostgres(t *testing.T) {
	ctx := context.Background()

	migrationsPath := findMigrationsPath(t)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "usdt",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
		),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresC.Terminate(ctx)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/usdt?sslmode=disable", host, port.Port())

	// pool
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// run migrations using your package
	require.NoError(t, postgres.Migrate(dsn, migrationsPath))

	repo := postgres.New(pool)

	rec := postgres.RateRecord{
		FetchedAt: time.Date(2026, 4, 9, 7, 0, 0, 0, time.UTC),
		AskTopN:   decimal.NewFromInt(80),
		AskAvgNM:  decimal.NewFromInt(81),
		BidTopN:   decimal.NewFromInt(79),
		BidAvgNM:  decimal.NewFromInt(78),
	}

	require.NoError(t, repo.SaveRate(ctx, rec))

	// verify columns exist and values match
	var got postgres.RateRecord
	row := pool.QueryRow(ctx, `
		SELECT fetched_at, ask_top_n, ask_avg_n_m, bid_top_n, bid_avg_n_m
		FROM rates
		ORDER BY id DESC
		LIMIT 1
	`)

	var (
		fetchedAt time.Time
		askTopN   decimal.Decimal
		askAvgNM  decimal.Decimal
		bidTopN   decimal.Decimal
		bidAvgNM  decimal.Decimal
	)

	// pgx scans into decimal.Decimal works with shopspring/decimal
	err = row.Scan(&fetchedAt, &askTopN, &askAvgNM, &bidTopN, &bidAvgNM)
	require.NoError(t, err)

	got.FetchedAt = fetchedAt
	got.AskTopN = askTopN
	got.AskAvgNM = askAvgNM
	got.BidTopN = bidTopN
	got.BidAvgNM = bidAvgNM

	require.True(t, got.FetchedAt.Equal(rec.FetchedAt))
	require.True(t, got.AskTopN.Equal(rec.AskTopN))
	require.True(t, got.AskAvgNM.Equal(rec.AskAvgNM))
	require.True(t, got.BidTopN.Equal(rec.BidTopN))
	require.True(t, got.BidAvgNM.Equal(rec.BidAvgNM))
}
