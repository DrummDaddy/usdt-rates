# USDT Rates gRPC Service (Grinex → PostgreSQL)

## Overview
Сервис предоставляет gRPC-метод **GetRates**, который:
1) запрашивает **спотовую depth** у биржи **Grinex** по HTTP (через `resty`)
2) вычисляет значения:
    - **topN** — значение на позиции `N` (1-indexed)
    - **avgNM** — среднее по диапазону `[N; M]` (включительно, 1-indexed)
3) сохраняет результат **в PostgreSQL** с меткой времени получения курса (каждый вызов `GetRates`)

Дополнительно:
- `Healthcheck` — проверка работоспособности
- graceful shutdown
- трассировка OpenTelemetry (через `otelgrpc` StatsHandler)
- логирование `zap`

---

## gRPC API

### Service
В `usdt.proto` сервис называется:

- `usdt.v1.proto.RateService`

Методы:
- `GetRates`
- `Healthcheck`

> Числа в ответе возвращаются **строками** (как в `proto`).

---

### GetRates
`GetRates(GetRatesRequest) returns (GetRatesResponse)`

#### Request 
```
{
  "n": 1,
  "m": 2
}
```
n — позиция (1-indexed) для topN и начало диапазона для avgNM
m — конец диапазона для avgNM (включительно, 1-indexed)
символ глубины в сервисе фиксирован: usdta7a5

### Response
```json{
  "fetched_at": "2026-04-09T07:00:00Z",
  "ask_top_n": "80.73",
  "ask_avg_n_m": "80.74",
  "bid_top_n": "79.73",
  "bid_avg_n_m": "79.72"
}
```
### Healthcheck
Healthcheck(HealthcheckRequest) returns (HealthcheckResponse)
```Response
json{ "status": "ok" }
```
### PostgreSQL
Схема создаётся миграциями из:

internal/storage/migrations

Таблица:

rates

Колонки:
```
fetched_at
ask_top_n, ask_avg_n_m
bid_top_n, bid_avg_n_m
```

### Конфигурация
Поддерживаются CLI flags и environment variables.
gRPC

--grpc-addr / GRPC_ADDR (по умолчанию :50051)

### PostgreSQL

--postgres-dsn / POSTGRES_DSN

### Grinex
````
--grinex-depth-url / GRINEX_DEPTH_URL
(по умолчанию: https://grinex.io/api/v1/spot/depth)
--grinex-timeout / GRINEX_TIMEOUT
(по умолчанию: 5s)
````
### OpenTelemetry (опционально)
````
--otel-service-name / OTEL_SERVICE_NAME
--otel-otlp-grpc-endpoint / OTEL_OTLP_GRPC_ENDPOINT
если OTEL_OTLP_GRPC_ENDPOINT пустой — tracing отключается
````

### Запуск
1) Unit-тесты
bashmake test
(Если make недоступен, используйте go test ./...)
2) Сборка
bashmake build
3) Запуск локально (нужен PostgreSQL)
bash./app

### Docker
Сборка образа 
````
bashdocker build -t usdt-rates:latest .
Запуск через docker-compose
Файл docker-compose.yml поднимает:

postgres
app

Команды:
bashdocker-compose up -d
docker-compose run --rm app ./app 
````