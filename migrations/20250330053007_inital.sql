-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM ('guest', 'staff','admin');
CREATE TYPE event_status AS ENUM('planned', 'ongoing', 'completed');
SELECT 'up SQL query';
CREATE TABLE users (
    telegram_id BIGINT PRIMARY KEY,
    firstname VARCHAR(30),
    lastname  varchar(30),
    patronymic varchar(30),
    chat_id BIGINT,
    role role_type
);
CREATE TABLE event(
    id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    speaker text,
    auditorium text,
    date timestamptz,
    status event_status DEFAULT 'planned'
);
CREATE TABLE user_event(
    telegram_id BIGINT REFERENCES users(telegram_id),
    schedule_id int REFERENCES event(id)
);
CREATE TABLE states(
    chat_id BIGINT PRIMARY KEY,
    data JSONB,
    state TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS user_event;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS event;
DROP TABLE IF EXISTS states;
DROP TYPE IF EXISTS role_type;
DROP TYPE IF EXISTS event_status;
-- +goose StatementEnd
