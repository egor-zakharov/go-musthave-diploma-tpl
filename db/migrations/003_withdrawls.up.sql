CREATE TABLE IF NOT EXISTS withdrawals
(
    id           BIGSERIAL PRIMARY KEY,
    number       VARCHAR references orders(number),
    sum          NUMERIC,
    user_id      bigint references users (id),
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);