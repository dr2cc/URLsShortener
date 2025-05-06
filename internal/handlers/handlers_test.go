package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func TestGetHandler(t *testing.T) {
	//Здесь общие для всех тестов данные
	shortURL := "6ba7b811"
	record := map[string]string{shortURL: "https://practicum.yandex.ru/"}

	tests := []struct {
		name       string
		method     string
		input      *storage.URLStorage
		want       string
		wantStatus int
	}{
		{
			name:   "all good",
			method: http.MethodGet,
			input: &storage.URLStorage{
				Data: record,
			},
			want:       "https://practicum.yandex.ru/",
			wantStatus: http.StatusTemporaryRedirect,
		},
		{
			name:   "with bad method",
			method: http.MethodPost,
			input: &storage.URLStorage{
				Data: record,
			},
			want:       "Method not allowed",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "key in input does not match /6ba7b811",
			method: http.MethodGet,
			input: &storage.URLStorage{
				Data: map[string]string{"6ba7b81": "https://practicum.yandex.ru/"},
			},
			want:       "URL with such id doesn't exist",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/"+shortURL, nil) //body)
			rr := httptest.NewRecorder()
			//***************************
			//Создаю экземпляр хранилища
			//05.05.2025- это не нужно!
			// работает только потому, что в реквесте не правильные данные
			// они пустые, вот и находят по пустому ключу!!
			//storage.MakeEntry(tt.input, "", "https://practicum.yandex.ru/")
			//!!!Вызов MakeEntry использовать для теста в storage
			//***************************
			handler := http.HandlerFunc(GetHandler(tt.input))

			handler.ServeHTTP(rr, req)

			if gotStatus := rr.Code; gotStatus != tt.wantStatus {
				t.Errorf("Want status '%d', got '%d'", tt.wantStatus, gotStatus)
			}

			// Ожидаемое (want) сообщение о ошибке должно совпадать с получаемым (got)
			if gotLocation := strings.TrimSpace(rr.Header()["Location"][0]); gotLocation != tt.want {
				t.Errorf("Want location'%s', got '%s'", tt.want, gotLocation)
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	// type args struct {
	// 	//w   http.ResponseWriter
	// 	w   *httptest.ResponseRecorder
	// 	req *http.Request
	// }

	//Здесь общие для всех тестов данные
	shortURL := "6ba7b811"
	record := map[string]string{shortURL: "https://practicum.yandex.ru/"}

	tests := []struct {
		name   string
		ts     *storage.URLStorage
		method string
		//args args
		//
		statusCode int
	}{
		{
			name: "all good",
			ts: &storage.URLStorage{
				Data: record,
			},
			method:     "POST",
			statusCode: http.StatusCreated,
		},
		{
			name: "bad method",
			ts: &storage.URLStorage{
				Data: record,
			},
			method: "GET",

			statusCode: http.StatusBadRequest,
		},

		// //Оставляю до сдачи инкремента 1
		// {
		// 	name: "invalid data type in query body- json",
		// 	ts: &storage.URLStorage{
		// 		Data: record,
		// 	},
		// 	method: "POST",
		// 	args: args{
		// 		w: httptest.NewRecorder(),
		// 		req: &http.Request{
		// 			Method: "POST",
		// 			Header: http.Header{
		// 				"Content-Type": []string{"applicatin/json"},
		// 			},
		// 			Host: host,
		// 			Body: body,
		// 		},
		// 	},
		// 	statusCode: http.StatusBadRequest,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// //Здесь использование
			// req, err := http.NewRequest(tt.method, host, body) //("POST", "/users/123", nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// // Позволяло пройти тесты,
			// // хотя host (со значением:= "localhost:8080") здесь это бред, он не давал нужных данных
			// // Данные видимо получались из
			// // body := io.NopCloser(bytes.NewBuffer([]byte(record[shortURL])))
			// // Но body спокойно заменилось nil , когда target:= "/"+shortURL)
			req := httptest.NewRequest(tt.method, "/"+shortURL, nil)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(PostHandler(tt.ts))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("Want status '%d', got '%d'", status, tt.statusCode)
			}
		})
	}
}
