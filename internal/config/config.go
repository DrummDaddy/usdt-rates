package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	GRPCAddr string

	PostgresDSN string

	GrinexDepthURL string
	GrinexTimeout  time.Duration

	OTelOTLPGRPCEndpoint string
	OTelServiceName      string
}

func New() Config {
	pflag.String("grpc-addr", ":50051", "gRPC listen address")

	pflag.String("postgres-dsn", "postgres://postgres:postgres@localhost:5432/usdt?sslmode=disable", "PostgreSQL DSN")

	pflag.String("grinex-depth-url", "https://grinex.io/api/v1/spot/depth", "Grinex spot depth endpoint")
	pflag.Duration("grinex-timeout", 5*time.Second, "HTTP timeout for Grinex requests")

	pflag.String("otel-service-name", "usdt-rates", "OpenTelemetry service name")
	pflag.String("otel-otlp-grpc-endpoint", "", "OTLP/gRPC endpoint, e.g. http://otel-collector:4317 (optional)")

	pflag.Parse()
	viper.AutomaticEnv()

	viper.BindEnv("grpc-addr", "GRPC_ADDR")
	viper.BindEnv("postgres-dsn", "POSTGRES_DSN")

	viper.BindEnv("grinex-depth-url", "GRINEX_DEPTH_URL")
	viper.BindEnv("grinex-timeout", "GRINEX_TIMEOUT")

	viper.BindEnv("otel-service-name", "OTEL_SERVICE_NAME")
	viper.BindEnv("otel-otlp-grpc-endpoint", "OTEL_OTLP_GRPC_ENDPOINT")

	viper.SetDefault("grpc-addr", ":50051")
	viper.SetDefault("postgres-dsn", "postgres://postgres:postgres@localhost:5432/usdt?sslmode=disable")
	viper.SetDefault("grinex-depth-url", "https://grinex.io/api/v1/spot/depth")
	viper.SetDefault("grinex-timeout", "5s")
	viper.SetDefault("otel-service-name", "usdt-rates")
	viper.SetDefault("otel-otlp-grpc-endpoint", "")

	return Config{
		GRPCAddr:    viper.GetString("grpc-addr"),
		PostgresDSN: viper.GetString("postgres-dsn"),

		GrinexDepthURL: viper.GetString("grinex-depth-url"),
		GrinexTimeout:  viper.GetDuration("grinex-timeout"),

		OTelOTLPGRPCEndpoint: viper.GetString("otel-otlp-grpc-endpoint"),
		OTelServiceName:      viper.GetString("otel-service-name"),
	}
}
