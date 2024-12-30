package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Dostonlv/url-shortner-http3-based/internal/models"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(url *models.URL) error {
	query := `
        INSERT INTO urls (short_code, original_url, created_at, clicks)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	return r.db.QueryRow(
		query,
		url.ShortCode,
		url.OriginalURL,
		url.CreatedAt,
		url.Clicks,
	).Scan(&url.ID)
}

func (r *URLRepository) Get(shortCode string) (*models.URL, error) {
	url := &models.URL{}
	query := `
        SELECT id, short_code, original_url, created_at, clicks
        FROM urls
        WHERE short_code = $1`

	err := r.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.Clicks,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No URL found for short code: %s", shortCode)
			return nil, fmt.Errorf("URL not found")
		}
		return nil, err
	}

	return url, nil
}

func (r *URLRepository) IncrementClicks(shortCode string) error {
	query := `
        UPDATE urls
        SET clicks = clicks + 1
        WHERE short_code = $1`

	_, err := r.db.Exec(query, shortCode)
	return err
}
