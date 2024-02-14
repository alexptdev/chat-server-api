-- +goose Up
-- +goose StatementBegin

create table chat_users (
    chat_users_chat_id    int,
    chat_users_user_id    int,
    chat_users_created_at timestamp not null default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_users cascade;
-- +goose StatementEnd
