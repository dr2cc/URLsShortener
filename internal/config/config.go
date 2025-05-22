package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// //**FromYandex
// // переменная FlagRunAddr содержит адрес и порт для запуска сервера
// var FlagRunAddr string

// // переменная FlagURL отвечает за базовый адрес результирующего сокращённого URL
// var FlagURL string

// // ParseFlags обрабатывает аргументы командной строки
// // и сохраняет их значения в соответствующих переменных
// func ParseFlags() {
// 	// регистрируем переменную flagRunAddr
// 	// как аргумент -a со значением :8080 по умолчанию
// 	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
// 	flag.StringVar(&FlagURL, "b", "http://localhost:8080", "host and port")

// 	// разбираем переданные серверу аргументы в зарегистрированные переменные
// 	flag.Parse()
// 	// Добавляю переменные окружения
// 	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
// 		//fmt.Println(envRunAddr)
// 		FlagRunAddr = envRunAddr
// 	}

// 	if envURL := os.Getenv("BASE_URL"); envURL != "" {
// 		//fmt.Println(envURL)
// 		FlagURL = envURL
// 	}
// }
// //**************

// анмаршаллить..
type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {

	// $env:CONFIG_PATH = "C:\__git\URLsShortener\config\local.yaml"       (на drkk)
	// $env:CONFIG_PATH = "C:\Mega\__git\URLsShortener\config\local.yaml"  (на ноуте)

	// // Если будут проблемы с переменной окружения, то писать путь так (\ экранируется \\):
	//configPath := "C:\\__git\\URLsShortener\\config\\local.yaml"

	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
