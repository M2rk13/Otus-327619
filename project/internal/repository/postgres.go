package repository

import (
	"context"
	"fmt"

	"currency-converter/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, dsn string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &PostgresRepository{db: pool}, nil
}

func (r *PostgresRepository) SaveConversion(ctx context.Context, history *model.ConversionHistory) error {
	query := `INSERT INTO conversion_history (currency_from, currency_to, amount, result, rate)
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, history.From, history.To, history.Amount, history.Result, history.Rate)
	if err != nil {
		return fmt.Errorf("failed to save conversion to postgres: %w", err)
	}
	return nil
}
