-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE SCHEMA ticket;

CREATE TABLE ticket.tickets (
    uuid UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    recipient_type VARCHAR(50) NOT NULL,
    recipient_uuid UUID NOT NULL,
    responsible_uuid UUID,
    priority INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE ticket.tickets;
DROP SCHEMA ticket;