package models

import (
	"time"

	"github.com/google/uuid"
)

// Модель транзакции, которая хранится в БД
type Transaction struct {
	ID          uuid.UUID
	FromAddress uuid.UUID
	ToAddress   uuid.UUID
	Amount      float32
	Status      Status
	Message     string
	CreatedAt   time.Time
}

// Модель для API-запроса на создание транзакции
type CreateTransactionRequest struct {
	FromAddress string  `json:"from" validate:"required,uuid"`
	ToAddress   string  `json:"to" validate:"required,uuid"`
	Amount      float32 `json:"amount" validate:"required,min=0"`
}

// Модель для ответа на API-запрос получения списка транзакций
type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	FromAddress uuid.UUID `json:"from"`
	ToAddress   uuid.UUID `json:"to"`
	Amount      float32   `json:"amount"`
	Status      Status    `json:"status"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}

// Модель аккумулирующая в себе параметры запроса для получения списка транзакций
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
		Message:     transaction.Message,
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
			Message:     transaction.Message,
			CreatedAt:   transaction.CreatedAt,
		}
		transactionResponses = append(transactionResponses, tr)
	}

	return transactionResponses
}

func ToTransaction(transaction *CreateTransactionRequest, id uuid.UUID, status Status, createdAt time.Time, message string) *Transaction {
	return &Transaction{
		ID:          id,
		FromAddress: uuid.MustParse(transaction.FromAddress),
		ToAddress:   uuid.MustParse(transaction.ToAddress),
		Amount:      transaction.Amount,
		Status:      status,
		Message:     message,
		CreatedAt:   createdAt,
	}
}
