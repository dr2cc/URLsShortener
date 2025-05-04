package storage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLStorage_GetHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      *URLStorage
		want       string
		statusCode int
	}{
		{
			name:   "all good",
			method: http.MethodGet,
			input: &URLStorage{
				Data: map[string]string{"6ba7b811": "https://practicum.yandex.ru/"},
			},
			want:       "https://practicum.yandex.ru/",
			statusCode: http.StatusTemporaryRedirect,
		},
		{
			name:   "with bad method",
			method: http.MethodPost,
			input: &URLStorage{
				Data: map[string]string{"6ba7b811": "https://practicum.yandex.ru/"},
			},
			want:       "Method not allowed",
			statusCode: http.StatusBadRequest,
		},
		// {
		// 	name:   "key in input does not match /6ba7b811",
		// 	method: http.MethodGet,
		// 	input: &URLStorage{
		// 		Data: map[string]string{"6ba7b81": "https://practicum.yandex.ru/"},
		// 	},
		// 	want:       "URL with such id doesn`t exist",
		// 	statusCode: http.StatusBadRequest,
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			// //Если запрос от клиента, то можно использовать пакет http. Не сработало..
			// request, _ := http.NewRequest(tc.method, "/6ba7b811", nil)
			request := httptest.NewRequest(tc.method, "/6ba7b811", nil)
			// Вызываем метод GetHandler структуры URLStorage (input)
			// Этот метод делает запись в responseRecorder
			tc.input.GetHandler(responseRecorder, request)

			// По заданию на конечную точку с методом GET в инкременте 1
			// в случае успешной обработки запроса сервер возвращает:

			// статус с кодом 307, должен совпадать с тем, что описан в statusCode
			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			// URL (переданный в input) в заголовке "Location", в случае ошибки,
			// сообщение о ошибке должно совпадать с want
			if strings.TrimSpace(responseRecorder.Header()["Location"][0]) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

// Основной каркас сгенерирован VSCode
// func TestURLStorage_PostHandler(t *testing.T) {
// 	type args struct {
// 		//w   http.ResponseWriter
// 		w   *httptest.ResponseRecorder
// 		req *http.Request
// 	}

// 	//Здесь стандартно передаваемые ("правильные") данные, вне теста получаемые от клиента
// 	host := "localhost:8080"
// 	shortURL := "6ba7b811"
// 	record := map[string]string{shortURL: "https://practicum.yandex.ru/"}
// 	body := io.NopCloser(bytes.NewBuffer([]byte(record[shortURL])))

// 	tests := []struct {
// 		name string
// 		ts   *URLStorage
// 		args args
// 		//
// 		statusCode int
// 	}{
// 		{
// 			name: "all good",
// 			ts: &URLStorage{
// 				Data: record,
// 			},
// 			args: args{
// 				w: httptest.NewRecorder(),
// 				//req: httptest.NewRequest("POST", "/6ba7b811", nil),
// 				req: &http.Request{
// 					Method: "POST",
// 					Header: http.Header{
// 						"Content-Type": []string{"text/plain"},
// 					},
// 					Host: host,
// 					Body: body,
// 				},
// 			},
// 			//если все нормально:
// 			//возвращает статус с кодом 201 (http.StatusCreated)
// 			statusCode: http.StatusCreated,
// 		},
// 		// {
// 		// 	name: "bad method",
// 		// 	ts: &URLStorage{
// 		// 		Data: record,
// 		// 	},
// 		// 	args: args{
// 		// 		w: httptest.NewRecorder(),
// 		// 		req: &http.Request{
// 		// 			Method: "GET",
// 		// 			Header: http.Header{
// 		// 				"Content-Type": []string{"text/plain"},
// 		// 			},
// 		// 			Host: host,
// 		// 			Body: body,
// 		// 		},
// 		// 	},
// 		// 	statusCode: http.StatusBadRequest,
// 		// },
// 		{
// 			name: "bad header",
// 			ts: &URLStorage{
// 				Data: record,
// 			},
// 			args: args{
// 				w: httptest.NewRecorder(),
// 				req: &http.Request{
// 					Method: "POST",
// 					Header: http.Header{
// 						"Content-Type": []string{"applicatin/json"},
// 					},
// 					Host: host,
// 					Body: body,
// 				},
// 			},
// 			statusCode: http.StatusBadRequest,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.ts.PostHandler(tt.args.w, tt.args.req)
// 			//Проверяю статус
// 			if tt.args.w.Code != tt.statusCode {
// 				t.Errorf("Want status '%d', got '%d'", tt.statusCode, tt.args.w.Code)
// 			}
// 			// //Здесь проверял содержимое body
// 			// //Когда починил рандомайзер, стала бессмысленной
// 			// if strings.TrimSpace(tt.args.w.Body.String()) != tt.want {
// 			// 	t.Errorf("Want '%s', got '%s'", tt.want, tt.args.w.Body)
// 			// }
// 		})
// 	}
// }
