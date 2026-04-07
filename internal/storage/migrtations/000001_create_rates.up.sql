CREATE TABLE IF NOT EXISTS rates (
    id BIGSERIAL PRIMARY KEY,
    fetched_at TIMESTAMP NOT NULL,

    ask_top_n NUMERIC(40, 18) NOT NULL,
    ask_avg_n_m NUMERIC(40, 18) NOT NULL,

    big_top_n NUMERIC(40, 18) NOT NULL,
    bid_avg_n_m NUMERIC(40, 18) NOT NULL
);

CREATE INDEX IF NOT EXISTS rates_fetched_at_idx ON rates (fetched_at);
