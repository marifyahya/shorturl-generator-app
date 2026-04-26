package service

import (
	"context"
	"testing"

	"github.com/marifyahya/shorturl-generator-app/internal/model"
)

// MockRepository is a manual mock implementation of repository.URLRepository
type MockRepository struct {
	CreateFunc         func(ctx context.Context, url *model.URL) error
	GetByShortCodeFunc func(ctx context.Context, code string) (*model.URL, error)
	IncrementHitsFunc  func(ctx context.Context, code string) error
}

func (m *MockRepository) Create(ctx context.Context, url *model.URL) error {
	return m.CreateFunc(ctx, url)
}

func (m *MockRepository) GetByShortCode(ctx context.Context, code string) (*model.URL, error) {
	return m.GetByShortCodeFunc(ctx, code)
}

func (m *MockRepository) IncrementHits(ctx context.Context, code string) error {
	return m.IncrementHitsFunc(ctx, code)
}

func TestURLService_Shorten_Success(t *testing.T) {
	mockRepo := &MockRepository{
		GetByShortCodeFunc: func(ctx context.Context, code string) (*model.URL, error) {
			return nil, nil // Not found, so code is available
		},
		CreateFunc: func(ctx context.Context, url *model.URL) error {
			return nil
		},
	}

	svc := NewURLService(mockRepo)
	code, err := svc.Shorten(context.Background(), "https://example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected code length 6, got %d", len(code))
	}
}

func TestURLService_Shorten_Collision(t *testing.T) {
	attempts := 0
	mockRepo := &MockRepository{
		GetByShortCodeFunc: func(ctx context.Context, code string) (*model.URL, error) {
			attempts++
			if attempts == 1 {
				return &model.URL{ShortCode: code}, nil // Found on first try (collision)
			}
			return nil, nil // Not found on second try
		},
		CreateFunc: func(ctx context.Context, url *model.URL) error {
			return nil
		},
	}

	svc := NewURLService(mockRepo)
	_, err := svc.Shorten(context.Background(), "https://example.com")

	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts due to collision, got %d", attempts)
	}
}

func TestURLService_GetOriginalURL_Success(t *testing.T) {
	expectedURL := "https://example.com"
	mockRepo := &MockRepository{
		GetByShortCodeFunc: func(ctx context.Context, code string) (*model.URL, error) {
			return &model.URL{OriginalURL: expectedURL}, nil
		},
		IncrementHitsFunc: func(ctx context.Context, code string) error {
			return nil
		},
	}

	svc := NewURLService(mockRepo)
	url, err := svc.GetOriginalURL(context.Background(), "abc123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url != expectedURL {
		t.Errorf("expected %s, got %s", expectedURL, url)
	}
}

func TestURLService_GetOriginalURL_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		GetByShortCodeFunc: func(ctx context.Context, code string) (*model.URL, error) {
			return nil, nil
		},
	}

	svc := NewURLService(mockRepo)
	_, err := svc.GetOriginalURL(context.Background(), "missing")

	if err == nil {
		t.Fatal("expected error for missing code, got nil")
	}
	if err.Error() != "short code not found" {
		t.Errorf("expected 'short code not found', got '%v'", err)
	}
}
