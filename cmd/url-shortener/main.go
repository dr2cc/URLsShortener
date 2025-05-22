package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dr2cc/URLsShortener.git/internal/config"
	"github.com/dr2cc/URLsShortener.git/internal/http-server/handlers/url/save"
	mwLogger "github.com/dr2cc/URLsShortener.git/internal/http-server/middleware/logger"
	"github.com/dr2cc/URLsShortener.git/internal/lib/logger/sl"
	"github.com/dr2cc/URLsShortener.git/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Перед запуском нужно установить переменную окружения CONFIG_PATH
//
// $env:CONFIG_PATH = "C:\__git\URLsShortener\config\local.yaml"
// $env:CONFIG_PATH = "C:\Mega\__git\URLsShortener\config\local.yaml"  (на ноуте)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	//
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	//
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	//переопределяем внутренний логгер
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	router.Post("/", save.New(log, storage))

	fmt.Println("Server is up!")

	// //**FromYandex
	// // мой "старый" код
	// // обрабатываем аргументы командной строки
	// config.ParseFlags()

	// if err := server.Run(); err != nil {
	// 	panic(err)
	// }
	// //**************
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
