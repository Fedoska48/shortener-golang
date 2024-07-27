package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "shortener-golang/internal/lib/api/response"
	"shortener-golang/internal/lib/logger/sl"
	"shortener-golang/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLDelete interface {
	DelURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requiest_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode req body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("req body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			log.Info("alias empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		deletedAlias, err := urlDelete.DelURL(alias)

		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("url not found", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("url", slog.String("alias", deletedAlias))

		render.JSON(w, r, Response{
			Response: resp.Delete(),
			Alias:    deletedAlias,
		})
	}
}
