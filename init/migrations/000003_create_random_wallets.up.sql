CREATE EXTENSION IF NOT EXISTS "pgcrypto";

TRUNCATE TABLE wallets RESTART IDENTITY CASCADE;

INSERT INTO wallets (id, balance)
SELECT gen_random_uuid()::VARCHAR(64), 100.00
FROM generate_series(1, 10);