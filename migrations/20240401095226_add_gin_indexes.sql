-- +goose Up
-- +goose StatementBegin
create index idx_users_first_name_tsvector on users using gin(first_name_tsvector);
create index idx_users_second_name_tsvector on users using gin(second_name_tsvector);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_users_first_name_tsvector;
drop index idx_users_second_name_tsvector;
-- +goose StatementEnd
