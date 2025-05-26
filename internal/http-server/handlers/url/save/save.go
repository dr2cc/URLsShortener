package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	// для краткости даем короткий алиас пакету
	resp "github.com/dr2cc/URLsShortener.git/internal/lib/api/response"
	"github.com/dr2cc/URLsShortener.git/internal/lib/logger/sl"
	"github.com/dr2cc/URLsShortener.git/internal/lib/random"
	"github.com/dr2cc/URLsShortener.git/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

// TODO: move to config if needed
const aliasLength = 6

// Объект urlSaver передадим при создании хендлера из main
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем объект запроса и десиариализуем (анмаршаллим) в него запрос
		//
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request")) // <----

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request")) // <----

			return
		}

		// Лучше больше логов, чем меньше - лишнее мы легко сможем почистить,
		// при необходимости. А вот недостающую информацию мы уже не получим.
		log.Info("request body decoded", slog.Any("req", req))

		// "Валидируем" запрос.
		// Нужно проверить, что URL — это действительно URL, что он не пустой.
		//
		// Создаем объект валидатора
		// и передаем в него структуру req (которую нужно провалидировать)
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// Короткий идентификатор (по которому будем искать оригинальный адрес),
		// проверяем вручную.
		// Если он пустой — генерируем случайный:
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		//  Сохраняем URL и Alias,
		//
		// Объект urlSaver (переданный при создании хендлера из main)
		// используется именно тут!
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// Отдельно обрабатываем ситуацию,
			// когда запись с таким Alias уже существует
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		// возвращаем ответ с сообщением об успехе
		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
