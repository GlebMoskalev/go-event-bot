-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    telegram_id BIGINT PRIMARY KEY,
    firstname VARCHAR(30),
    lastname  varchar(30),
    patronymic varchar(30),
    is_admin bool
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user_schedule;
DROP TABLE users;
DROP TABLE schedule;
-- +goose StatementEnd
