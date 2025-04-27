package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func TestPostHandler(t *testing.T) {
	storage := storage.NewStorage()
	handler := PostHandler(storage)

	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		wantStatus  int
	}{
		{
			name:        "Valid POST",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        "https://example.com",
			wantStatus:  http.StatusCreated,
		},
		{
			name:        "Invalid content type",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        "https://example.com",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "Invalid method",
			method:      http.MethodGet,
			contentType: "text/plain",
			body:        "https://example.com",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	storage := storage.NewStorage()
	handler := GetHandler(storage)

	// Setup test data
	testID := "test123"
	testURL := "https://example.com"
	storage.InsertURL(testID, testURL)

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "Valid GET",
			method:       http.MethodGet,
			path:         "/" + testID,
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: testURL,
		},
		{
			name:         "Non-existent ID",
			method:       http.MethodGet,
			path:         "/nonexistent",
			wantStatus:   http.StatusBadRequest,
			wantLocation: "URL with such id doesn't exist",
		},
		{
			name:         "Invalid method",
			method:       http.MethodPost,
			path:         "/" + testID,
			wantStatus:   http.StatusBadRequest,
			wantLocation: "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			location := resp.Header.Get("Location")
			if location != tt.wantLocation {
				t.Errorf("Expected Location %s, got %s", tt.wantLocation, location)
			}
		})
	}
}

func TestGenerateShortURL(t *testing.T) {
	storage := storage.NewStorage()
	url := "https://example.com/very/long/url"

	short := generateShortURL(storage, url)
	if len(short) < 2 {
		t.Errorf("Generated URL is too short: %s", short)
	}

	// Verify the URL was stored
	id := strings.TrimPrefix(short, "/")
	if _, err := storage.GetURL(id); err != nil {
		t.Errorf("URL was not stored in storage: %v", err)
	}
}
