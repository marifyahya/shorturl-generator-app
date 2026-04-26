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

func (h *URLHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, errorResponse{Error: message})
}

func (h *URLHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
