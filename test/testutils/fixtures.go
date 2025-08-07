package testutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура загрузчика SQL миграций
// Используется для загрузки и накатывания миграций-фикстур на testContainer
// Хранит в себе базовый путь до тестовых файлов и pgx connection pool
// 
// Используется только для тестирования
type FixtureManager struct {
	pool     *pgxpool.Pool
	basePath string
}

func NewFixtureManager(pool *pgxpool.Pool) *FixtureManager {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return &FixtureManager{
		pool: pool,
		basePath: filepath.Join(dir, "..", "testdata"),
	}
}

//Функция для загрузки и накатывания фикстур
func (fm *FixtureManager) ApplySQLFixture(ctx context.Context, fixturePath string) error {
	fullPath := filepath.Join(fm.basePath, "fixtures", "sql", fixturePath)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL fixture: %w", err)
	}

	commands := strings.Split(string(content), ":")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == ""{
			continue
		}
		if _, err := fm.pool.Exec(ctx, cmd); err != nil {
			return fmt.Errorf("failed to execute SQL command %s: %w", cmd, err)
		}
	}
	return nil
}