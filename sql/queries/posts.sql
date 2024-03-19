-- name: CreatePost :one
insert into posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostsByUserID :many
select * from posts
where feed_id in (select id from feeds where user_id = $1)
order by published_at desc limit $2;