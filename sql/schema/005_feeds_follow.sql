-- +goose Up
create table feed_follow (
    id uuid primary key,
    feed_id uuid not null references feeds(id) on delete cascade on update cascade ,
    user_id uuid not null references users(id) on delete cascade on update cascade ,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table feed_follow;