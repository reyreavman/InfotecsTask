package postgres

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/pkg/database"
)

type TransactionRepository struct {
	db *database.Client
}

func NewTransactionRepository(db *database.Client) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r TransactionRepository) GetTransactions(ctx context.Context, count int) ([]*models.Transaction, error) {
	sql := `SELECT id, from_address, to_address, amount, status, created_at 
            FROM transactions 
            ORDER BY created_at DESC 
            LIMIT $1`

	rows, err := r.db.Query(ctx, sql, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0, count)
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.FromAddress, &t.ToAddress, &t.Amount, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r TransactionRepository) GetAllTransactions(ctx context.Context) ([]*models.Transaction, error) {
	sql := `SELECT id, from_address, to_address, amount, status, created_at 
            FROM transactions 
            ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.FromAddress, &t.ToAddress, &t.Amount, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}