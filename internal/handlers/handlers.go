package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func generateShortURL(urlList *storage.URLStorage, longURL string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := []rune(longURL)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	//[:11] здесь сокращаю строку
	id := reg.ReplaceAllString(string(runes[:11]), "")

	storage.MakeEntry(urlList, id, longURL)

	return "/" + id
}

// Функция PostHandler уровня пакета handlers
func PostHandler(ts *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		param, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Преобразуем тело запроса (тип []byte) в строку:
		longURL := string(param)
		// Генерируем сокращённый URL и создаем запись в нашем хранилище
		shortURL := "http://" + req.Host + generateShortURL(ts, longURL)

		// Устанавливаем статус ответа 201
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, shortURL)
	}
}

// Функция GetHandler уровня пакета handlers
func GetHandler(ts *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := strings.TrimPrefix(req.RequestURI, "/")
		longURL, err := storage.GetEntry(ts, id)
		if err != nil {
			w.Header().Set("Location", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
