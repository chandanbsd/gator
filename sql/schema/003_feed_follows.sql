-- +goose Up
CREATE TABLE feed_follows(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);

ALTER TABLE feed_follows
ADD CONSTRAINT unique_user_feed unique (user_id, feed_id);

-- +goose Down
DROP TABLE feed_follows;
