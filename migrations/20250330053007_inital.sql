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
CREATE TABLE schedule (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    description text,
    date timestamptz
);
CREATE TABLE user_schedule(
    telegram_id BIGINT REFERENCES users(telegram_id),
    schedule_id int REFERENCES schedule(id)
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
DROP TABLE user_schedule;
DROP TABLE users;
DROP TABLE schedule;
DROP TABLE staffs;
DROP TYPE role_type;
-- +goose StatementEnd
