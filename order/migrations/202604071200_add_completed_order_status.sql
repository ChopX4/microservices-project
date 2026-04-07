-- +goose Up
ALTER TYPE order_status ADD VALUE IF NOT EXISTS 'COMPLETED';

-- +goose Down

