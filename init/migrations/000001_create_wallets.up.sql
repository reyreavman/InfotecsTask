CREATE TABLE wallets
(
    id          VARCHAR(64) PRIMARY KEY,
    balance     DECIMAL(10, 2) NOT NULL DEFAULT 100.00
);

COMMENT ON TABLE wallets IS 'Таблица для хранения информации о кошельках';
COMMENT ON COLUMN wallets.id IS 'Идентификатор кошелька';
COMMENT ON COLUMN wallets.balance IS 'Баланс кошелька';