package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dr2cc/URLsShortener.git/internal/config"
	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func generateAlias(us *storage.URLStorage, url string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := []rune(url)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	//[:11] здесь сокращаю строку
	id := reg.ReplaceAllString(string(runes[:11]), "")

	storage.MakeEntry(us, id, url)

	return "/" + id
}

// Функция PostHandler уровня пакета handlers
func PostHandler(us *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			param, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Преобразуем тело запроса (тип []byte) в строку:
			url := string(param)
			// // Генерируем короткий идентификатор и создаем запись в нашем хранилище
			//alias := "http://" + req.Host + generateAlias(us, url)

			alias := config.FlagURL + generateAlias(us, url)

			// Устанавливаем статус ответа 201
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, alias)

		default:
			w.Header().Set("Location", "Method not allowed")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// Функция GetHandler уровня пакета handlers
func GetHandler(us *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			id := strings.TrimPrefix(req.RequestURI, "/")
			url, err := storage.GetEntry(us, id)
			if err != nil {
				w.Header().Set("Location", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.Header().Set("Location", "Method not allowed")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
