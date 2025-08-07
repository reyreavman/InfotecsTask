package postgres

import (
	"context"
	"errors"
	"fmt"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/pkg/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Реализация репозитория
type PaymentRepository struct {
	db *database.Client
}

func NewPaymentRepository(db *database.Client) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// Реализация метода для создания транзакций
//
// Если не найден кошелёк отправителя или получателя возвращается ошибка и запись в БД не создается
// Если кошельки найдены, в БД создается запись о транзакции со статусом pending и соответствующим сообщением
//
// В случае, когда на балансе отправителя не хватает нужной суммы для совершения транзакции,
// запись о транзакции в БД обновляется со статусом failed и соответствующим сообщением
//
// В случае, если все необходимые условия выполнены,
// запись о транзакции в БД обновляется со статусом completed и соответствующим сообщением
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
			if errors.Is(err, pgx.ErrNoRows) {
				return payment.ErrRecipientWalletNotFound
			}
			return fmt.Errorf("failed to lock recipient wallet: %w", err)
		}

		transaction = &models.Transaction{
			ID:          uuid.New(),
			FromAddress: uuid.MustParse(createTransactionRequest.FromAddress),
			ToAddress:   uuid.MustParse(createTransactionRequest.ToAddress),
			Amount:      createTransactionRequest.Amount,
			Status:      models.Pending,
			Message:     "Transaction Pending",
			CreatedAt:   time.Now(),
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO transactions (id, from_address, to_address, amount, status, message, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			transaction.ID,
			transaction.FromAddress,
			transaction.ToAddress,
			transaction.Amount,
			transaction.Status,
			transaction.Message,
			transaction.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		if senderBalance < createTransactionRequest.Amount {
			transaction.Status = models.Failed
			transaction.Message = models.SENDER_NOT_HAVE_ENOUGH_BALANCE
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1, message = $2 WHERE id = $3`,
				transaction.Status,
				transaction.Message,
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
			transaction.Message = models.TRANSACTION_FAILED
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1, message = $2 WHERE id = $3`,
				transaction.Status,
				transaction.Message,
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
			transaction.Message = models.TRANSACTION_FAILED
			_, err = tx.Exec(
				ctx,
				`UPDATE transactions SET status = $1, message = $2 WHERE id = $3`,
				transaction.Status,
				transaction.Message,
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
		transaction.Message = models.TRANSACTION_COMPLETED
		_, err = tx.Exec(
			ctx,
			`UPDATE transactions SET status = $1, message = $2 WHERE id = $3`,
			transaction.Status,
			transaction.Message,
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
