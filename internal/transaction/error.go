package transaction

import "errors"

var ErrNonPositiveTransactionAmount = errors.New("Transaction amount should be positive")

var ErrFailedToCreateTransaction = errors.New("Failed to create transaction")