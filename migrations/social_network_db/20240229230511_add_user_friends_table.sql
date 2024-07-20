-- +goose Up
-- +goose StatementBegin
create table user_friends
(
    user_id text,
    friend_user_id text,
    primary key (user_id, friend_user_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friends
-- +goose StatementEnd
