package repository

import "github.com/Dostonlv/url-shortner-http3-based/internal/models"

type URLRepository interface {
	Create(url *models.URL) error
	Get(shortCode string) (*models.URL, error)
	IncrementClicks(shortCode string) error
}
