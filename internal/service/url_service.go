package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/marifyahya/shorturl-generator-app/internal/model"
	"github.com/marifyahya/shorturl-generator-app/internal/repository"
	"strings"
)

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, code string) (string, error)
	GetStats(ctx context.Context, code string) (*model.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(ctx context.Context, originalURL string) (string, error) {
	originalURL = strings.TrimSpace(originalURL)
	if originalURL == "" {
		return "", errors.New("url cannot be empty")
	}

	// Auto-fix: Add https:// if no protocol is provided
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		code := GenerateShortCode()

		// Check if code already exists
		existing, err := s.repo.GetByShortCode(ctx, code)
		if err != nil {
			return "", fmt.Errorf("error checking existing code: %v", err)
		}

		if existing == nil {
			// Code is available
			url := &model.URL{
				ShortCode:   code,
				OriginalURL: originalURL,
			}
			if err := s.repo.Create(ctx, url); err != nil {
				return "", fmt.Errorf("error creating url: %v", err)
			}
			return code, nil
		}
	}

	return "", errors.New("could not generate a unique short code after several attempts")
}

func (s *urlService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	url, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("error getting original url: %v", err)
	}
	if url == nil {
		return "", errors.New("short code not found")
	}

	// Increment hits (best effort or check error)
	_ = s.repo.IncrementHits(ctx, code)

	return url.OriginalURL, nil
}

func (s *urlService) GetStats(ctx context.Context, code string) (*model.URL, error) {
	url, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %v", err)
	}
	if url == nil {
		return nil, errors.New("short code not found")
	}
	return url, nil
}
