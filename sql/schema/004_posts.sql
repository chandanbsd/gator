-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	description TEXT NULL,
	published_at TIMESTAMP NULL,
	feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);

ALTER TABLE posts
ADD CONSTRAINT unique_posts_url unique (url);

-- +goose Down
DROP TABLE posts;
