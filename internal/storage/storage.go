package storage

import (
	"errors"
	"net/http"
	"strings"
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

// !!!Попробовать пройти автотесты с этим методом!!
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
