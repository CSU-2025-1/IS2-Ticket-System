-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS notifications;

CREATE TABLE IF NOT EXISTS notifications.mail_receiver
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mail VARCHAR(512) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications.mail_receiver;
-- +goose StatementEnd
