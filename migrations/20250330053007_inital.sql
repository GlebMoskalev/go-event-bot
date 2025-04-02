-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM ('guest', 'staff','admin');
SELECT 'up SQL query';
CREATE TABLE users (
    telegram_id BIGINT PRIMARY KEY,
    firstname VARCHAR(30),
    lastname  varchar(30),
    patronymic varchar(30),
    role role_type
);
CREATE TABLE event(
    id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    speaker text,
    auditorium text,
    date timestamptz
);
CREATE TABLE user_event(
    telegram_id BIGINT REFERENCES users(telegram_id),
    schedule_id int REFERENCES event(id)
);
CREATE TABLE staffs (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(30),
    lastname  varchar(30),
    patronymic varchar(30),
    email varchar(50),
    phone_number varchar(20)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user_event;
DROP TABLE users;
DROP TABLE event;
DROP TABLE staffs;
DROP TYPE role_type;
-- +goose StatementEnd
