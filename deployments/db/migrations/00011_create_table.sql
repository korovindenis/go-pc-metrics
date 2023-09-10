-- +goose Up
-- SQL in Up.
-- Description: This migration creates the gauge and counter tables.
CREATE TABLE gauge (
    id SERIAL PRIMARY KEY,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name CHAR(50) NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

CREATE TABLE counter (
    id SERIAL PRIMARY KEY,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name CHAR(50) NOT NULL,
    delta INTEGER NOT NULL
);

-- +goose Down
-- SQL in Down.
-- Description: This migration drops the gauge and counter tables.
DROP TABLE counter;
DROP TABLE gauge;