package db

import (
	"context"
	"fmt"
	"time"

	"albums-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB — пул соединений с PostgreSQL
type DB struct {
	Pool *pgxpool.Pool
}

// NewConnection создаёт новое подключение к базе
func NewConnection(cfg *config.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать DSN: %w", err)
	}

	// Настройки пула соединений
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	// Проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не удалось пингануть базу: %w", err)
	}

	fmt.Println("✅ Успешно подключились к PostgreSQL")
	return &DB{Pool: pool}, nil
}

// Close закрывает пул соединений
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		fmt.Println("🔌 Соединение с PostgreSQL закрыто")
	}
}
