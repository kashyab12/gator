-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetNextFeedsToFetch :many
select * from feeds order by last_fetched_at limit $1;

-- name: MarkFeedFetched :one
update feeds
set last_fetched_at = $1,
    updated_at = $1
where id = $2
returning *;