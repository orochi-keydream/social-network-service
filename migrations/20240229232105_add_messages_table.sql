-- +goose Up
-- +goose StatementBegin
create table messages
(
    message_id bigserial primary key,
    sent_at timestamp not null,
    from_user_id text not null,
    to_user_id text not null,
    text text not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages
-- +goose StatementEnd
