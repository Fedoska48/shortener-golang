package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "shortener-golang/internal/lib/api/response"
	"shortener-golang/internal/lib/logger/sl"
	"shortener-golang/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requiest_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)

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

		log.Info("got url by alias", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
