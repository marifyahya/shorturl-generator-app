package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/marifyahya/shorturl-generator-app/internal/config"
	"github.com/marifyahya/shorturl-generator-app/internal/service"
)

type URLHandler struct {
	svc service.URLService
	cfg *config.Config
}

func NewURLHandler(svc service.URLService, cfg *config.Config) *URLHandler {
	return &URLHandler{
		svc: svc,
		cfg: cfg,
	}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Shorten handles the creation of a short URL.
// @Summary Shorten a long URL
// @Description Take a long URL and return a unique 6-character short code
// @Accept  json
// @Produce  json
// @Param   request body shortenRequest true "URL to shorten"
// @Success 201 {object} shortenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/shorten [post]
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	longURL := strings.TrimSpace(req.URL)
	if longURL == "" {
		h.respondWithError(w, http.StatusBadRequest, "url is required")
		return
	}

	// Validate URL format
	u, err := url.ParseRequestURI(longURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		h.respondWithError(w, http.StatusBadRequest, "invalid url format")
		return
	}

	code, err := h.svc.Shorten(r.Context(), longURL)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shortURL := h.cfg.BaseURL + code
	h.respondWithJSON(w, http.StatusCreated, shortenResponse{ShortURL: shortURL})
}

// Redirect handles the redirection from a short code to the original URL.
// @Summary Redirect to original URL
// @Description Take a 6-character short code and redirect the user to the original long URL
// @Param   short_code path string true "Short Code"
// @Success 302 {string} string "Found"
// @Failure 404 {object} errorResponse
// @Router /{short_code} [get]
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Simple path parsing for /:short_code
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		h.respondWithError(w, http.StatusBadRequest, "short code is required")
		return
	}

	originalURL, err := h.svc.GetOriginalURL(r.Context(), code)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "short code not found")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// GetStats handles the retrieval of statistics for a short URL.
// @Summary Get statistics for a short URL
// @Description Return metadata and hit count for a given short code
// @Produce  json
// @Param   short_code path string true "Short Code"
// @Success 200 {object} model.URL
// @Failure 404 {object} errorResponse
// @Router /api/stats/{short_code} [get]
func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Simple path parsing for /api/stats/:short_code
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		h.respondWithError(w, http.StatusBadRequest, "short code is required")
		return
	}
	code := parts[3]

	urlData, err := h.svc.GetStats(r.Context(), code)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "short code not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, urlData)
}

func (h *URLHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, errorResponse{Error: message})
}

func (h *URLHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
