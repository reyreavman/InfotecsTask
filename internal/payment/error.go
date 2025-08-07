package payment

import "errors"

//Список ошибок бизнес-логики при которых выполнение транзакции завершается, если вызвана одна из этих ошибок, запись транзакции в БД сохранена не будет
var ErrSenderWalletNotFound = errors.New("Sender wallet not found")
var ErrRecipientWalletNotFound = errors.New("Recipient wallet not found")