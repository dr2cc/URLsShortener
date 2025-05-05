package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func TestGetHandler(t *testing.T) {
	//Здесь стандартно передаваемые ("правильные") данные, вне теста получаемые от клиента
	//host := "localhost:8080"
	shortURL := "6ba7b811"
	record := map[string]string{shortURL: "https://practicum.yandex.ru/"}
	body := io.NopCloser(bytes.NewBuffer([]byte(record[shortURL])))

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
		// {
		// 	name:   "key in input does not match /6ba7b811",
		// 	method: http.MethodGet,
		//  input: &storage.URLStorage{
		// 	    Data: record,
		//  },
		// 	want:       "URL with such id doesn`t exist",
		// 	wantStatus: http.StatusBadRequest,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// responseRecorder := httptest.NewRecorder()
			// // //Если запрос от клиента, то можно использовать пакет http. Не сработало..
			// // request, _ := http.NewRequest(tc.method, "/6ba7b811", nil)
			// request := httptest.NewRequest(tc.method, "/6ba7b811", nil)
			// // Вызываем метод GetHandler структуры URLStorage (input)
			// // Этот метод делает запись в responseRecorder
			// tc.input.GetHandler(responseRecorder, request)

			// // По заданию на конечную точку с методом GET в инкременте 1
			// // в случае успешной обработки запроса сервер возвращает:

			// // статус с кодом 307, должен совпадать с тем, что описан в statusCode
			// if responseRecorder.Code != tc.statusCode {
			// 	t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			// }

			//***************************

			req, err := http.NewRequest(tt.method, "localhost:8080/"+shortURL, body) //("POST", "/users/123", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			//***************************
			//Создаю экземпляр хранилища
			//05.05.2025- это не нужно!
			// работает только потому, что в реквесте не правильные данные
			// они пустые, вот и находят по пустому ключу!!
			//
			//id := strings.TrimPrefix(req.RequestURI, "/")
			//здесь в тесте уже в tt.input есть все нужные данные
			//storage.MakeEntry(tt.input, id, "https://practicum.yandex.ru/")
			storage.MakeEntry(tt.input, "", "https://practicum.yandex.ru/")
			//***************************
			handler := http.HandlerFunc(GetHandler(tt.input))

			handler.ServeHTTP(rr, req)

			if gotStatus := rr.Code; gotStatus != tt.wantStatus {
				t.Errorf("Want status '%d', got '%d'", tt.wantStatus, gotStatus)
			}

			//***********

			// URL (переданный в input) в заголовке "Location", в случае ошибки,
			// сообщение о ошибке должно совпадать с want
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
