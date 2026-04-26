package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marifyahya/shorturl-generator-app/internal/config"
	"github.com/marifyahya/shorturl-generator-app/internal/model"
)

// MockService is a manual mock for service.URLService
type MockService struct {
	ShortenFunc        func(ctx context.Context, originalURL string) (string, error)
	GetOriginalURLFunc func(ctx context.Context, code string) (string, error)
	GetStatsFunc       func(ctx context.Context, code string) (*model.URL, error)
}

func (m *MockService) Shorten(ctx context.Context, originalURL string) (string, error) {
	return m.ShortenFunc(ctx, originalURL)
}

func (m *MockService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	return m.GetOriginalURLFunc(ctx, code)
}

func (m *MockService) GetStats(ctx context.Context, code string) (*model.URL, error) {
	return m.GetStatsFunc(ctx, code)
}

func TestURLHandler_Shorten(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080/"}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockShorten    func(ctx context.Context, originalURL string) (string, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: shortenRequest{URL: "https://google.com"},
			mockShorten: func(ctx context.Context, originalURL string) (string, error) {
				return "abc123", nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"short_url":"http://localhost:8080/abc123"}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
		{
			name:           "Empty URL",
			requestBody:    shortenRequest{URL: ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"url is required"}`,
		},
		{
			name:           "Invalid URL Format",
			requestBody:    shortenRequest{URL: "not-a-url"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid url format"}`,
		},
		{
			name:        "Service Error",
			requestBody: shortenRequest{URL: "https://google.com"},
			mockShorten: func(ctx context.Context, originalURL string) (string, error) {
				return "", errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"db error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{ShortenFunc: tt.mockShorten}
			h := NewURLHandler(mockSvc, cfg)

			var body []byte
			if s, ok := tt.requestBody.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			h.Shorten(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestURLHandler_Redirect(t *testing.T) {
	tests := []struct {
		name           string
		shortCode      string
		mockGetURL     func(ctx context.Context, code string) (string, error)
		expectedStatus int
		expectedLocation string
	}{
		{
			name:      "Success",
			shortCode: "abc123",
			mockGetURL: func(ctx context.Context, code string) (string, error) {
				return "https://google.com", nil
			},
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://google.com",
		},
		{
			name:      "Not Found",
			shortCode: "missing",
			mockGetURL: func(ctx context.Context, code string) (string, error) {
				return "", errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{GetOriginalURLFunc: tt.mockGetURL}
			h := NewURLHandler(mockSvc, nil)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.shortCode, nil)
			rr := httptest.NewRecorder()

			h.Redirect(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusFound {
				location := rr.Header().Get("Location")
				if location != tt.expectedLocation {
					t.Errorf("expected location %s, got %s", tt.expectedLocation, location)
				}
			}
		})
	}
}
