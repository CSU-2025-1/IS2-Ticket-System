-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
create schema if not exists auth;

create table if not exists auth.users_auth_data
(
    uuid     uuid unique,
    login    varchar(64),
    password varchar(512)
);

insert into auth.users_auth_data(uuid, login, password)
values ('', 'admin', 'ed5faf26d7f8370b05e5df423baa316060090dfbf1cea4e30f8a2ebeaad22e268e7f157e46d5e2f3f89e5de146142175e051b9791744fb338fc69241cfa04006');

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop schema if exists auth;
drop table if exists auth.users_auth_data;