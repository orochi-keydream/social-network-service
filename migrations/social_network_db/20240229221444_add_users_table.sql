-- +goose Up
-- +goose StatementBegin
create table users
(
    user_id text primary key,
    first_name text not null,
    second_name text not null,
    gender int not null,
    birthdate date not null,
    biography text not null,
    city text not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users
-- +goose StatementEnd
