package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

// 04.05.25 автотесты прошли!
// Закончить здесь, а затем перенести в handler
// func TestGetHandler(t *testing.T) {
// 	tt := []struct {
// 		name       string
// 		method     string
// 		input      *storage.URLStorage
// 		want       string
// 		statusCode int
// 	}{
// 		{
// 			name:   "all good",
// 			method: http.MethodGet,
// 			input: &storage.URLStorage{
// 				Data: map[string]string{"6ba7b811": "https://practicum.yandex.ru/"},
// 			},
// 			want:       "https://practicum.yandex.ru/",
// 			statusCode: http.StatusTemporaryRedirect,
// 		},
// 		{
// 			name:   "with bad method",
// 			method: http.MethodPost,
// 			input: &storage.URLStorage{
// 				Data: map[string]string{"6ba7b811": "https://practicum.yandex.ru/"},
// 			},
// 			want:       "Method not allowed",
// 			statusCode: http.StatusBadRequest,
// 		},
// 		// {
// 		// 	name:   "key in input does not match /6ba7b811",
// 		// 	method: http.MethodGet,
// 		// 	input: &URLStorage{
// 		// 		Data: map[string]string{"6ba7b81": "https://practicum.yandex.ru/"},
// 		// 	},
// 		// 	want:       "URL with such id doesn`t exist",
// 		// 	statusCode: http.StatusBadRequest,
// 		// },
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			responseRecorder := httptest.NewRecorder()
// 			// //Если запрос от клиента, то можно использовать пакет http. Не сработало..
// 			// request, _ := http.NewRequest(tc.method, "/6ba7b811", nil)
// 			request := httptest.NewRequest(tc.method, "/6ba7b811", nil)
// 			// Вызываем метод GetHandler структуры URLStorage (input)
// 			// Этот метод делает запись в responseRecorder
// 			tc.input.GetHandler(responseRecorder, request)

// 			// По заданию на конечную точку с методом GET в инкременте 1
// 			// в случае успешной обработки запроса сервер возвращает:

// 			// статус с кодом 307, должен совпадать с тем, что описан в statusCode
// 			if responseRecorder.Code != tc.statusCode {
// 				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
// 			}

// 			// URL (переданный в input) в заголовке "Location", в случае ошибки,
// 			// сообщение о ошибке должно совпадать с want
// 			if strings.TrimSpace(responseRecorder.Header()["Location"][0]) != tc.want {
// 				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
// 			}
// 		})
// 	}
// }

func TestPostHandler(t *testing.T) {
	// type args struct {
	// 	//w   http.ResponseWriter
	// 	w   *httptest.ResponseRecorder
	// 	req *http.Request
	// }

	//Здесь стандартно передаваемые ("правильные") данные, вне теста получаемые от клиента
	host := "localhost:8080"
	shortURL := "6ba7b811"
	record := map[string]string{shortURL: "https://practicum.yandex.ru/"}
	body := io.NopCloser(bytes.NewBuffer([]byte(record[shortURL])))

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
		// 	name: "bad header",
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
			req, err := http.NewRequest(tt.method, host, body) //("POST", "/users/123", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(PostHandler(tt.ts))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("Want status '%d', got '%d'", status, tt.statusCode)
			}
		})
	}
}
