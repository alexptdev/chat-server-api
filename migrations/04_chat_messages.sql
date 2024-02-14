-- +goose Up
-- +goose StatementBegin

create table chat_messages (
    message_id         serial primary key,
    message_chat_id    int,
    message_from       varchar(35),
    message_text       varchar(1024),
    message_created_at timestamp not null default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_messages cascade;
-- +goose StatementEnd
