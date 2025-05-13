package main

import (
	"github.com/dr2cc/URLsShortener.git/internal/config"
	"github.com/dr2cc/URLsShortener.git/internal/server"
)

func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
