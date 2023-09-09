-- +goose Up
CREATE TABLE gauge (
    id SERIAL PRIMARY KEY,
    name CHAR(20) UNIQUE,
    value DOUBLE PRECISION
);

CREATE TABLE counter (
    id SERIAL PRIMARY KEY,
    name CHAR(20) UNIQUE,
    value INTEGER
);

-- +goose Down
DROP TABLE counter;
DROP TABLE gauge;