package handlers

import (
	"net/http"
	"reflect"
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

func TestPostHandler(t *testing.T) {
	type args struct {
		ts *storage.URLStorage
	}

	// //Здесь стандартно передаваемые ("правильные") данные, вне теста получаемые от клиента
	// host := "localhost:8080"
	shortURL := "6ba7b811"
	// responseRecorder := httptest.NewRecorder()
	storageInstance := storage.NewStorage()
	storage.MakeEntry(storageInstance, shortURL, "https://practicum.yandex.ru/")
	//record := map[string]string{shortURL: "https://practicum.yandex.ru/"}
	// body := io.NopCloser(bytes.NewBuffer([]byte(record[shortURL])))

	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "помогите, спасите",
			args: args{
				ts: *storageInstance,
			},
			// args: &storage.URLStorage{
			// 	map[string]string{shortURL: "https://practicum.yandex.ru/"}
			// },
			want: func(w http.ResponseWriter, r *http.Request) {
				// w: &,
				// r: &http.Request{
				// 	Method: "POST",
				// 	Header: http.Header{
				// 		"Content-Type": []string{"applicatin/json"},
				// 	},
				// 	Host: host,
				// 	Body: body,
				// },
			},
			//statusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PostHandler(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
