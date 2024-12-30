package service

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Dostonlv/url-shortner-http3-based/internal/models"
	"github.com/Dostonlv/url-shortner-http3-based/internal/repository"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// generateShortCode generatsiya funksiyasi
func generateShortCode() string {
	const (
		alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		length   = 6
	)

	result := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			// Xato bo'lsa, vaqt asosida generatsiya qilamiz
			return time.Now().Format("150405")
		}
		result[i] = alphabet[n.Int64()]
	}

	return string(result)
}

func (s *URLService) CreateShortURL(originalURL string) (*models.URL, error) {
	url := &models.URL{
		ShortCode:   generateShortCode(),
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetURL(shortCode string) (*models.URL, error) {
	return s.repo.Get(shortCode)
}

func (s *URLService) IncrementClicks(shortCode string) error {
	return s.repo.IncrementClicks(shortCode)
}
