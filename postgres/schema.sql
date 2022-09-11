CREATE TABLE orders
(
    id SERIAL PRIMARY KEY,
    collective VARCHAR(50) NOT NULL,
    guild_id VARCHAR(50) NOT NULL, -- server_id
    created_at TIMESTAMPTZ NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'NTD',
    description TEXT NOT NULL
);

CREATE TABLE order_participants
(
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id varchar(50) NOT NULL,
    price MONEY NOT NULL -- a user should paid basic
);