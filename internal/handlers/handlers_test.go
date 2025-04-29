package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

// func TestGetHandler(t *testing.T) {
// 	type args struct {
// 		ts *storage.URLStorage
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want http.HandlerFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetHandler(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetHandler() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGetHandler(t *testing.T) {
// 	type args struct {
// 		ts *storage.URLStorage
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		want       http.HandlerFunc
// 		urlID      string
// 		wantStatus int
// 		wantLoc    string
// 	}{
// 		{
// 			name: "existing URL",
// 			args: args{func() *storage.URLStorage {
// 				s := storage.NewStorage()
// 				s.InsertURL("abc123", "https://example.com")
// 				return s
// 			}()},
// 			want: func() http.HandlerFunc {
// 				return func(w http.ResponseWriter, r *http.Request) {}
// 			}(),
// 			urlID:      "abc123",
// 			wantStatus: http.StatusTemporaryRedirect,
// 			wantLoc:    "https://example.com",
// 		},
// 		{
// 			name: "non-existing URL",
// 			args: args{ts: storage.NewStorage()},
// 			want: func() http.HandlerFunc {
// 				return func(w http.ResponseWriter, r *http.Request) {}
// 			}(),
// 			urlID:      "nonexistent",
// 			wantStatus: http.StatusBadRequest,
// 			wantLoc:    "URL with such id doesn't exist",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handler := GetHandler(tt.args.ts)

// 			req := httptest.NewRequest("GET", "/"+tt.urlID, nil)
// 			w := httptest.NewRecorder()

// 			handler(w, req)

// 			resp := w.Result()
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("GetHandler() status = %v, want %v", resp.StatusCode, tt.wantStatus)
// 			}

// 			loc := w.Header().Get("Location")
// 			if loc != tt.wantLoc {
// 				t.Errorf("GetHandler() Location header = %v, want %v", loc, tt.wantLoc)
// 			}
// 		})
// 	}
// }

//*************************************************************************

func TestPostHandler(t *testing.T) {
	type args struct {
		ts *storage.URLStorage
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
		//
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name: "successful URL shortening",
			args: args{ts: storage.NewStorage()},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {}
			}(),
			body:       "https://example.com",
			wantStatus: http.StatusCreated,
			wantBody:   "http://localhost:8080/", // Будет дополнено ID
		},
		{
			name: "empty body",
			args: args{ts: storage.NewStorage()},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {}
			}(),
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantBody:   "EOF", // Ожидаемая ошибка при пустом теле
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := PostHandler(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("PostHandler() = %v, want %v", got, tt.want)
			// }
			//
			handler := PostHandler(tt.args.ts)

			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("PostHandler() status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			body := w.Body.String()
			if tt.wantStatus == http.StatusCreated {
				if !strings.HasPrefix(body, tt.wantBody) {
					t.Errorf("PostHandler() body = %v, want prefix %v", body, tt.wantBody)
				}
			} else {
				if !strings.Contains(body, tt.wantBody) {
					t.Errorf("PostHandler() body = %v, want contains %v", body, tt.wantBody)
				}
			}
		})
	}
}
