-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
create extension if not exists "uuid-ossp";

create schema if not exists users;
create table if not exists users.users(
    uuid uuid primary key default uuid_generate_v4(),
    login varchar(64)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop schema if exists users;
drop table if exists users.users;