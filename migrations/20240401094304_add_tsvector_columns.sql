-- +goose Up
-- +goose StatementBegin
alter table users
add first_name_tsvector tsvector,
add second_name_tsvector tsvector;

update users
set first_name_tsvector = to_tsvector('english', first_name),
    second_name_tsvector = to_tsvector('english', second_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
drop first_name_tsvector,
drop second_name_tsvector
-- +goose StatementEnd
