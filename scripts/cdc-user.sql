INSERT INTO users (email, display_name, metadata)
VALUES ('eve@example.com', 'Eve Evans', JSON_OBJECT('plan', 'free'));

UPDATE users SET login_count = login_count + 1, last_login = NOW() WHERE email = 'alice@example.com';

UPDATE users SET status = 'suspended' WHERE login_count < 5;

DELETE FROM users WHERE email = 'dave@example.com';

-- Ensure the CDC event captures this within the transaction GTID
START TRANSACTION;
INSERT INTO users (email, display_name) VALUES ('frank@example.com', 'Frank Foster');
UPDATE users SET status = 'active' WHERE email = 'carol@example.com';
COMMIT;
