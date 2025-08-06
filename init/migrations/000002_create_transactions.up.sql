CREATE TABLE transactions
(
    id              VARCHAR(64) PRIMARY KEY,
    from_address    VARCHAR(64) NOT NULL,
    to_address      VARCHAR(64) NOT NULL,
    amount          DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    status          VARCHAR NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (from_address) REFERENCES wallets(id),
    FOREIGN KEY (to_address) REFERENCES wallets(id)
);

CREATE INDEX tr_created_at_idx ON transactions (created_at DESC);

COMMENT ON TABLE transactions IS 'Таблица для хранения информации о транзакциях';
COMMENT ON COLUMN transactions.id IS 'Идентификатор транзакции';
COMMENT ON COLUMN transactions.from_address IS 'Идентификатор отправителя';
COMMENT ON COLUMN transactions.to_address IS 'Идентификатор получателя';
COMMENT ON COLUMN transactions.amount IS 'Сумма транзакции';
COMMENT ON COLUMN transactions.status IS 'Статус транзакции';
COMMENT ON COLUMN transactions.created_at IS 'Время создания транзакции';