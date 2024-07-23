-- +goose Up
-- +goose StatementBegin
create table posts
(
    post_id text primary key,
    published_at timestamp not null,
    user_id text not null,
    text text not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table posts
-- +goose StatementEnd
