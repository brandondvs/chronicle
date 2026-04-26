-- Seed schema for Chronicle CDC demo
-- This runs automatically on first MySQL container boot via /docker-entrypoint-initdb.d

USE chronicle_demo;

-- A representative OLTP table: a variety of column types to exercise the CDC pipeline.
-- Includes nullable columns, a JSON column, timestamps, and a soft-delete flag so you can
-- generate INSERT, UPDATE, and DELETE events that look like real application data.
CREATE TABLE users (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    email        VARCHAR(255)    NOT NULL,
    display_name VARCHAR(100)    NOT NULL,
    status       ENUM('active', 'suspended', 'deleted') NOT NULL DEFAULT 'active',
    metadata     JSON                NULL,
    login_count  INT UNSIGNED    NOT NULL DEFAULT 0,
    last_login   DATETIME            NULL,
    created_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_email (email),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Seed a few rows so you have something to UPDATE and DELETE against immediately.
INSERT INTO users (email, display_name, status, metadata, login_count, last_login) VALUES
    ('alice@example.com',   'Alice Anderson', 'active',    JSON_OBJECT('plan', 'pro',  'beta', true),  42, '2026-04-25 14:30:00'),
    ('bob@example.com',     'Bob Brown',      'active',    JSON_OBJECT('plan', 'free', 'beta', false),  3, '2026-04-20 09:15:00'),
    ('carol@example.com',   'Carol Chen',     'suspended', JSON_OBJECT('plan', 'pro',  'reason', 'payment_failed'), 128, '2026-04-10 22:45:00'),
    ('dave@example.com',    'Dave Davis',     'active',    NULL, 0, NULL);

-- Grant the chronicle user the replication privileges it needs to act as a fake replica
-- and read the binlog. Without these, syncer.StartSync() will fail with an access error.
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'chronicle'@'%';
GRANT SELECT ON chronicle_demo.* TO 'chronicle'@'%';
FLUSH PRIVILEGES;
