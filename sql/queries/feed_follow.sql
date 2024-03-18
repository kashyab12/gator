-- name: CreateFeedFollow :one
insert into feed_follow (id, feed_id, user_id, created_at, updated_at)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetFeedFollows :many
select * from feed_follow where user_id = $1;