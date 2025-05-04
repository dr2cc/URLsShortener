package storage

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Storager interface {
	InsertURL(uid string, url string) error
	GetURL(uid string) (string, error)
}

type URLStorage struct {
	Data map[string]string
}

func NewStorage() *URLStorage {
	return &URLStorage{
		Data: make(map[string]string),
	}
}

func (s *URLStorage) InsertURL(uid string, url string) error {
	s.Data[uid] = url
	return nil
}

// метод GetURL типа *URLStorage
func (s *URLStorage) GetURL(uid string) (string, error) {
	e, exists := s.Data[uid]
	if !exists {
		return uid, errors.New("URL with such id doesn't exist")
	}
	return e, nil
}

func generateShortURL(urlList *URLStorage, longURL string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := []rune(longURL)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	//[:11] здесь сокращаю строку
	id := reg.ReplaceAllString(string(runes[:11]), "")

	MakeEntry(urlList, id, longURL)

	return "/" + id
}

// тип *URLStorage и его метод PostHandler
func (s *URLStorage) PostHandler(w http.ResponseWriter, req *http.Request) {
	// Автотесты не проходили на еще одном уровне switch
	//не знаю как на этом, но без него проходит любой тип контента,
	//а возвратиться может только text
	switch req.Header.Get("Content-Type") {
	case "text/plain":
		//param - тело запроса (тип []byte)
		param, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Генерирую ответ и создаю запись в хранилище
		response := "http://" + req.Host + generateShortURL(s, string(param))

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, response)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Content-Type isn`t text/plain")
	}
}

// !!!Попробовать пройти автотесты с этим методом!!
// 04.05.25  прошел автотесты!
// метод GetHandler типа *URLStorage
// Получается, что так логичнее если нет пакета handlers !!
// А если есть, то правильнее функция уровня пакета (?!?)
func (s *URLStorage) GetHandler(w http.ResponseWriter, req *http.Request) {
	//Тесты подсказали добавить проверку на метод:
	switch req.Method {
	case http.MethodGet:
		// //Пока (14.04.2025) не знаю как передать PathValue при тестировании.
		// id := req.PathValue("id")

		// А вот RequestURI получается и от клиента и из теста
		// Но получаю лишний "/"
		id := strings.TrimPrefix(req.RequestURI, "/")

		//Реализую интерфейс
		longURL, err := GetEntry(s, id)

		if err != nil {
			//http.Error(w, "URL not found", http.StatusBadRequest)
			w.Header().Set("Location", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", longURL)
		// //И так и так работает. Оставил первоначальный вариант.
		//http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.Header().Set("Location", "Method not allowed")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Реализую интерфейс Storager
func MakeEntry(s Storager, uid string, url string) {
	s.InsertURL(uid, url)
}

func GetEntry(s Storager, uid string) (string, error) {
	return s.GetURL(uid)
}
