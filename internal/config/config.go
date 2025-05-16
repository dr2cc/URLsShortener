package config

import (
	"flag"
	"os"
)

// переменная FlagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string

// переменная FlagURL отвечает за базовый адрес результирующего сокращённого URL
var FlagURL string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagURL, "b", "http://localhost:8080", "host and port")

	// разбираем переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Добавляю переменные окружения
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envURL := os.Getenv("BASE_URL"); envURL != "" {
		FlagURL = envURL
	}
}
