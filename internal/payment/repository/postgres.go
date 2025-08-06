package postgres

import (
	"context"
	"errors"
	"fmt"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/pkg/database"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type PaymentRepository struct {
	db *database.Client
}

func NewPaymentRepository(db *database.Client) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, createTransactionRequest *models.CreateTransactionRequest) (*models.TransactionResponse, error) {
	var transaction *models.Transaction

	err := r.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var senderBalance, recipientBalance float32

		err := tx.QueryRow(
			ctx,
			`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`,
			createTransactionRequest.FromAddress,
		).Scan(&senderBalance)
		if err != nil {
			log.Printf("Ошибка при блокировке кошелька sender: %v", err)
			if errors.Is(err, pgx.ErrNoRows) {
				return payment.ErrSenderWalletNotFound
			}
			return fmt.Errorf("failed to lock sender wallet: %w", err)
		}

		err = tx.QueryRow(
			ctx,
			`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`,
			createTransactionRequest.ToAddress,
		).Scan(&recipientBalance)
		if err != nil {
			log.Printf("Ошибка при блокировке кошелька recipient: %v", err)
			if errors.Is(err, pgx.ErrNoRows) {
				return payment.ErrRecipientWalletNotFound
			}
			return fmt.Errorf("failed to lock recipient wallet: %w", err)
		}

		transaction = &models.Transaction{
			ID:          uuid.New(),
			FromAddress: createTransactionRequest.FromAddress,
			ToAddress:   createTransactionRequest.ToAddress,
			Amount:      createTransactionRequest.Amount,
			Status:      models.Pending,
			CreatedAt:   time.Now(),
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO transactions (id, from_address, to_address, amount, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			transaction.ID,
			transaction.FromAddress,
			transaction.ToAddress,
			transaction.Amount,
			transaction.Status,
			transaction.CreatedAt,
		)
		if err != nil {
			log.Printf("Ошибка при создании транзакции: %v", err)
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		if senderBalance < createTransactionRequest.Amount {
			transaction.Status = models.Failed
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1 WHERE id = $2`,
				transaction.Status,
				transaction.ID,
			)
			if err != nil {
				return fmt.Errorf("failed to update transaction: %w", err)
			}
			return nil
		}

		_, err = tx.Exec(
			ctx,
			`UPDATE wallets SET balance = balance - $1 WHERE id = $2`,
			createTransactionRequest.Amount,
			createTransactionRequest.FromAddress,
		)
		if err != nil {
			transaction.Status = models.Failed
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1 WHERE id = $2`,
				transaction.Status,
				transaction.ID,
			)
			return nil
		}

		_, err = tx.Exec(
			ctx,
			`UPDATE wallets SET balance = balance + $1 WHERE id = $2`,
			createTransactionRequest.Amount,
			createTransactionRequest.ToAddress,
		)
		if err != nil {
			transaction.Status = models.Failed
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1 WHERE id = $2`,
				transaction.Status,
				transaction.ID,
			)
			_, err = tx.Exec(
				ctx,
				`UPDATE wallets SET balance = balance + $1 WHERE id = $2`,
				createTransactionRequest.Amount,
				createTransactionRequest.FromAddress,
			)
			return nil
		}

		transaction.Status = models.Completed
		_, err = tx.Exec(
			ctx,
			`UPDATE transactions SET status = $1 WHERE id = $2`,
			transaction.Status,
			transaction.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return models.ToTransactionResponse(transaction), nil
}
