-- +goose Up
CREATE TYPE event_status AS ENUM ('PENDING', 'PUBLISHED');

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_uuid UUID NOT NULL UNIQUE,
    topic TEXT NOT NULL,
    key TEXT NOT NULL,
    payload BYTEA NOT NULL,
    status event_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE events;
DROP TYPE event_status;
