package storage

import (
	"errors"
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

// Реализую интерфейс Storager
func MakeEntry(s Storager, uid string, url string) {
	s.InsertURL(uid, url)
}

func GetEntry(s Storager, uid string) (string, error) {
	return s.GetURL(uid)
}
