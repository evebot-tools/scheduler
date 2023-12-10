package internal

import (
	"fmt"
	"github.com/evebot-tools/utils"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

func startApi() {
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	r := chi.NewRouter()
	apiServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", utils.GetEnv("PORT", "3000")),
		Handler: r,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(sentryMiddleware.Handle)

	r.Get("/{group}/{endpoint}", func(w http.ResponseWriter, r *http.Request) {
		err := scheduler.RunByTag(fmt.Sprintf("%s_%s", chi.URLParam(r, "group"), chi.URLParam(r, "endpoint")))
		if err != nil {
			log.Err(err).Msg("failed to run status job by tag")
			_, _ = w.Write([]byte(
				fmt.Sprintf(
					"failed to triggered sync of %s/%s endpoint",
					chi.URLParam(r, "group"),
					chi.URLParam(r, "endpoint"),
				),
			))
			return
		}
		_, _ = w.Write(
			[]byte(
				fmt.Sprintf(
					"triggered sync of %s/%s endpoint",
					chi.URLParam(r, "group"),
					chi.URLParam(r, "endpoint"),
				),
			),
		)
	})

	log.Info().Msg("starting webserver")
	err := apiServer.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start webserver")
	}

}
