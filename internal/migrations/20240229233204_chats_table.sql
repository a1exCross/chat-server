-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    usernames TEXT[]
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE chats;
DROP DATABASE messages;
-- +goose StatementEnd
