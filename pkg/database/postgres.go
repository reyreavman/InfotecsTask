package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Обвязка для postgres клиента
// Используется pgx connection pool
type Client struct {
	pool *pgxpool.Pool
}

func NewClientWithPool(pool *pgxpool.Pool) *Client {
	return &Client{
		pool: pool,
	}
}

func NewClient(ctx context.Context, config Config) (*Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Client{
		pool: pool,
	}, nil
}

func (db *Client) Close() {
	db.pool.Close()
}

func (db *Client) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Обвязка для pgx функции Exec
func (db *Client) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := db.pool.Exec(ctx, sql, args...)
	return err
}

// Обвязка для pgx функции Query
func (db *Client) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

// Обвязка для pgx функции QueryRow
func (db *Client) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

// Обвязка для pgx функции BeginTx
func (db *Client) ExecuteTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				log.Printf("rollback error during panic: %v", rollbackErr)
			}
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
