package server

import (
	"log/slog"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes configures all HTTP routes for the server.
func setupRoutes(router *chi.Mux, level slog.Level, checker healthChecker) {
	router.Use(logMiddleware(level))

	router.Get("/", indexHandler())
	router.Handle("/metrics", promhttp.Handler())
	router.Get("/health", healthHandler(checker))
	router.Get("/ready", readyHandler(checker))

	if level == slog.LevelDebug {
		slog.Warn("serving pprof path via '/debug/pprof'")
		router.Route("/debug/pprof", func(router chi.Router) {
			router.HandleFunc("/*", pprof.Index)
			router.HandleFunc("/cmdline", pprof.Cmdline)
			router.HandleFunc("/profile", pprof.Profile)
			router.HandleFunc("/symbol", pprof.Symbol)
			router.HandleFunc("/trace", pprof.Trace)
		})
	}
}
