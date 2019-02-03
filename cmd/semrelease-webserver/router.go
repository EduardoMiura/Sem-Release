package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/poc-git/semrelease"
	"go.uber.org/zap"
)

func createServerHandler(service *semrelease.Service, logger *zap.Logger) http.Handler {

	router := chi.NewRouter()

	router.Use(
		middleware.RealIP,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Timeout(10*time.Second),
		// transactionIDMiddleware,
		// setAccessControl,
	)

	router.Route("/semrelease", func(router chi.Router) {
		router.Get("/repository", errorWrapper(getAllRepositoryHandler(service), logger))
		router.Post("/webhook", errorWrapper(webhookHandler(service), logger))
	})

	router.Route("/config", func(router chi.Router) {
		router.Post("/add", errorWrapper(addConfigHandler(service), logger))
		router.Get("/config", errorWrapper(getAllConfigHandler(service), logger))
	})

	router.Get("/version", versionHandler)
	router.Get("/health", errorWrapper(healthCheckHandler(service), logger))

	return router

}

func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w,
		http.StatusText(http.StatusNotImplemented),
		http.StatusNotImplemented,
	)
}
