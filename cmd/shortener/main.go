package main

import (
	"github.com/dr2cc/URLsShortener.git/internal/config"
	"github.com/dr2cc/URLsShortener.git/internal/server"
)

// const (
// 	envLocal = "local"
// 	envDev   = "dev"
// 	envProd  = "prod"
// )

func main() {
	// //из примера
	// cfg := config.MustLoad()

	// //
	// log := setupLogger(cfg.Env)
	// log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	// log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	// log.Debug("logger debug mode enabled")

	// обрабатываем аргументы командной строки
	config.ParseFlags()

	if err := server.Run(); err != nil {
		panic(err)
	}
}

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// 	switch env {
// 	case envLocal:
// 		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
// 	case envDev:
// 		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
// 	case envProd:
// 		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
// 	}

// 	return log
// }
