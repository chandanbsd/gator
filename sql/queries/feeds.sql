-- name: CreateFeed :one
INSERT INTO feeds(
	id,
	created_at,
	updated_at,
	name,
	url,
	user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING *;

-- name: GetFeeds :many
select f.name, f.url, u.name as user_name
from feeds f
join users u on u.id = f.user_id;

-- name: GetFeedByUrl :one
select f.ID, f.url
from feeds f
where f.url = $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $1
where ID = $2;

-- name: GetNextFeedToFetch :one
select *
from feeds
order by last_fetched_at desc;
