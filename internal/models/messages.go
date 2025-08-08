package models

// Возможные сообщения внутри транзакций
const (
	SENDER_NOT_HAVE_ENOUGH_BALANCE = "Sender does not have enough balance"
	TRANSACTION_COMPLETED          = "Transaction completed"
	TRANSACTION_PENDING            = "Transaction pending"
	TRANSACTION_FAILED             = "Transaction failed"
)
