-- name: CreatePost :exec
INSERT INTO posts(
    id,
	created_at,
	updated_at,
	title,
	url,
	description,
	published_at,
	feed_id)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetPostsForUser :many
SELECT *
FROM posts p
WHERE EXISTS (
    SELECT 1
    FROM feed_follows ff
    WHERE ff.feed_id = p.feed_id
        AND ff.user_id = $1
)
ORDER BY published_at desc
LIMIT $2;
