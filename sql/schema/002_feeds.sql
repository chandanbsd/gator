-- +goose Up
CREATE TABLE feeds(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	name TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	last_fetched_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE feeds;
