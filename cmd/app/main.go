package main

import (
	"context"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DrummDaddy/usdt-rates/internal/config"
	appgrpc "github.com/DrummDaddy/usdt-rates/internal/grpc"
	"github.com/DrummDaddy/usdt-rates/internal/rates/client"
	"github.com/DrummDaddy/usdt-rates/internal/service"
	"github.com/DrummDaddy/usdt-rates/internal/storage/postgres"

	usdtpb "github.com/DrummDaddy/usdt-rates/gen/usdt/v1"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func initOTel(cfg config.Config, logger *zap.Logger) func() {
	if strings.TrimSpace(cfg.OTelOTLPGRPCEndpoint) == "" {
		return func() {}
	}
	endpoint := cfg.OTelOTLPGRPCEndpoint
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		u, err := url.Parse(endpoint)
		if err == nil && u.Host != "" {
			endpoint = u.Host
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		logger.Warn("Failed to initialize OTel", zap.Error(err))
		return func() {}
	}
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", cfg.OTelServiceName),
	),
	)
	if err != nil {
		logger.Warn("Failed to create OTel resource, tracing disabled", zap.Error(err))
		return func() {}
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = tp.Shutdown(shutdownCtx)
	}

}

func main() {
	cfg := config.New()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)

	}
	defer func() { _ = logger.Sync() }()

	shutdownOTel := initOTel(cfg, logger)
	defer shutdownOTel()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("Failed to initialize PostgreSQL pool", zap.Error(err))
	}
	defer pool.Close()

	const migrationPath = "internal/storage/migrations"
	if err := postgres.Migrate(cfg.PostgresDSN, migrationPath); err != nil {
		logger.Fatal("Failed to migrate", zap.Error(err))
	}

	grinex := client.NewGrinexClient(cfg.GrinexDepthURL, cfg.GrinexTimeout)

	repo := postgres.New(pool)
	svc := service.New(grinex, repo)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	handler := appgrpc.NewHandler(svc)

	usdtpb.RegisterRateServiceServer(grpcServer, handler)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting gRPC server", zap.String("addr", cfg.GRPCAddr))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()
	<-stopCh
	logger.Info("shutdown gRPC server")
	grpcServer.GracefulStop()

}
