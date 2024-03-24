-- +goose Up
-- +goose StatementBegin
create table user_accounts
(
    user_id text primary key,
    password_hash text not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_accounts
-- +goose StatementEnd
