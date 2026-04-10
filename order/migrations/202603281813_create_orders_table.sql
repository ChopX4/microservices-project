-- +goose Up
CREATE TYPE payment_method AS ENUM ('UNKNOWN', 'CARD', 'SPB', 'CREDIT_CARD', 'INVESTOR_MONEY');
CREATE TYPE order_status AS ENUM ('PENDING_PAYMENT', 'PAID', 'CANCELED');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_uuid UUID NOT NULL UNIQUE,
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price NUMERIC(12, 2) NOT NULL,
    transaction_uuid UUID NOT NULL,
    payment_method payment_method NOT NULL DEFAULT 'UNKNOWN',
    status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE orders;
DROP TYPE payment_method;
DROP TYPE order_status;