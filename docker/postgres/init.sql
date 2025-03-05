-- схема под Ory/Hydra
CREATE SCHEMA IF NOT EXISTS hydra;

-- Отдельная схема под каждый микросервис
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS gateway;
CREATE SCHEMA IF NOT EXISTS notification;
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS tickets;