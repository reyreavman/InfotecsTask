package testutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Структура, аккумулирующая в себе testContainer и pgx connection pool
// Имеет методы для запуска testContainer, его закрытия и накатывания миграций
type PGTestContainer struct {
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
}

func StartPGContainer(ctx context.Context, migrationsPath string) (*PGTestContainer, error) {
	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := applyMigrations(ctx, pool, migrationsPath); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &PGTestContainer{
		Container: container,
		Pool:      pool,
	}, nil
}

func (pg *PGTestContainer) Close(ctx context.Context) error {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
	if pg.Container != nil {
		return pg.Container.Terminate(ctx)
	}
	return nil
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, path string) error {
	log.Printf("Applying migrations from: %s", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			log.Printf("Migration content:\n%s", string(content))
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		log.Printf("Applied migration: %s", file.Name())
	}

	return nil
}
