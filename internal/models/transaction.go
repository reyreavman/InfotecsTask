package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID
	FromAddress uuid.UUID
	ToAddress   uuid.UUID
	Amount      float32
	Status      Status
	CreatedAt   time.Time
}

// Модель для API-запроса на создание транзакции (DTO)
type CreateTransactionRequest struct {
	FromAddress uuid.UUID `json:"from" validate:"required,uuid"`
	ToAddress   uuid.UUID `json:"to" validate:"required,uuid"`
	Amount      float32   `json:"amount" validate:"required,min=0"`
}

// Модель для ответа на API-запрос получения списка транзакций (DTO)
type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	FromAddress uuid.UUID `json:"from"`
	ToAddress   uuid.UUID `json:"to"`
	Amount      float32   `json:"amount"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetTransactionWithCountRequest struct {
	Count *int `form:"count" validate:"omitempty,min=1"`
}

func ToTransactionResponse(transaction *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:          transaction.ID,
		FromAddress: transaction.FromAddress,
		ToAddress:   transaction.ToAddress,
		Amount:      transaction.Amount,
		Status:      transaction.Status,
		CreatedAt:   transaction.CreatedAt,
	}
}

func ToTransactionResponses(transactions []*Transaction) []*TransactionResponse {
	transactionResponses := make([]*TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		tr := &TransactionResponse{
			ID:          transaction.ID,
			FromAddress: transaction.FromAddress,
			ToAddress:   transaction.ToAddress,
			Amount:      transaction.Amount,
			Status:      transaction.Status,
			CreatedAt:   transaction.CreatedAt,
		}
		transactionResponses = append(transactionResponses, tr)
	}

	return transactionResponses
}

func ToTransaction(transaction *CreateTransactionRequest, uuid uuid.UUID, status Status, createdAt time.Time) *Transaction {
	return &Transaction{
		ID:          uuid,
		FromAddress: transaction.FromAddress,
		ToAddress:   transaction.ToAddress,
		Amount:      transaction.Amount,
		Status:      status,
		CreatedAt:   createdAt,
	}
}
