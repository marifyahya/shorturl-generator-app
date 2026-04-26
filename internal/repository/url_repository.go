package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marifyahya/shorturl-generator-app/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	GetByShortCode(ctx context.Context, code string) (*model.URL, error)
	IncrementHits(ctx context.Context, code string) error
}

type postgresURLRepository struct {
	db *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) URLRepository {
	return &postgresURLRepository{db: db}
}

func (r *postgresURLRepository) Create(ctx context.Context, url *model.URL) error {
	query := `INSERT INTO urls (short_code, original_url) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, url.ShortCode, url.OriginalURL).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		return fmt.Errorf("could not create url: %v", err)
	}
	return nil
}

func (r *postgresURLRepository) GetByShortCode(ctx context.Context, code string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, hits, created_at FROM urls WHERE short_code = $1`
	var url model.URL
	err := r.db.QueryRowContext(ctx, query, code).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.Hits, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or return a custom error like ErrNotFound
		}
		return nil, fmt.Errorf("could not get url: %v", err)
	}
	return &url, nil
}

func (r *postgresURLRepository) IncrementHits(ctx context.Context, code string) error {
	query := `UPDATE urls SET hits = hits + 1 WHERE short_code = $1`
	_, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("could not increment hits: %v", err)
	}
	return nil
}
