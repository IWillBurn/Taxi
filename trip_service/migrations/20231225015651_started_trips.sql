-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS started_trips
(
    id INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    offer_id VARCHAR(255) NOT NULL,
    driver_id VARCHAR(255) NOT NULL,
    current_stage VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS started_trips
-- +goose StatementEnd