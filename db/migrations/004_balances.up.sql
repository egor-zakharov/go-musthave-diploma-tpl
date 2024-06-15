CREATE TABLE IF NOT EXISTS balances
(
    user_id bigint references users (id) UNIQUE,
    sum     NUMERIC
);