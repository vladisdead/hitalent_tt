-- +goose Up
-- +goose StatementBegin
create table if not exists chats (
    id serial primary key,
    tittle varchar(200) not null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);;
create table if not exists messages (
    id serial,
    chat_id int REFERENCES chats(id) on delete CASCADE,
    text varchar(5000) not null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists messages;
drop table if exists chats;
-- +goose StatementEnd
