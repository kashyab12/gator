-- +goose Up
create table feeds (
    id uuid primary key,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    name text not null,
    url text not null unique,
    user_id uuid not null,
    foreign key (user_id) references users(id)
        on delete cascade
        on update cascade
);

-- +goose Down
drop table feeds;
