-- +goose Up
-- +goose StatementBegin

create table chats (
    chat_id          serial primary key,
    chat_name        varchar(25),
    chat_description varchar(256),
    chat_author_id   int,
    chat_created_at  timestamp not null default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chats cascade;
-- +goose StatementEnd
