package storage

import (
	"errors"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

// // **FromYandex
// type Storager interface {
// 	InsertURL(id string, url string) error
// 	GetURL(id string) (string, error)
// }

// type URLStorage struct {
// 	Data map[string]string
// }

// func NewStorage() *URLStorage {
// 	return &URLStorage{
// 		Data: make(map[string]string),
// 	}
// }

// func (s *URLStorage) InsertURL(id string, url string) error {
// 	s.Data[id] = url
// 	return nil
// }

// // метод GetURL типа *URLStorage
// func (s *URLStorage) GetURL(id string) (string, error) {
// 	e, exists := s.Data[id]
// 	if !exists {
// 		return id, errors.New("URL with such id doesn't exist")
// 	}
// 	return e, nil
// }

// // Реализую интерфейс Storager
// func MakeEntry(s Storager, id string, url string) {
// 	s.InsertURL(id, url)
// }

// func GetEntry(s Storager, id string) (string, error) {
// 	return s.GetURL(id)
// }
