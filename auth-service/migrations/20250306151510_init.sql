-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
create schema if not exists auth;

create table if not exists auth.users_auth_data
(
    uuid     uuid unique,
    login    varchar(63),
    password varchar(255)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop schema if exists auth;
drop table if exists auth.users_auth_data;