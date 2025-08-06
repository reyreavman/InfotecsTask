package payment

import "errors"

var ErrSenderWalletNotFound = errors.New("Sender wallet not found")
var ErrRecipientWalletNotFound = errors.New("Recipient wallet not found")
var ErrSenderNotHaveEnoughBalance = errors.New("Sender does not have enough balance")