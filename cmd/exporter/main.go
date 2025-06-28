package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dirsigler/hackerone-exporter/internal/config"
	"github.com/dirsigler/hackerone-exporter/internal/exporter"
	"github.com/dirsigler/hackerone-exporter/internal/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "hackerone-exporter",
		Usage: "Export HackerOne metrics to Prometheus",
		Flags: config.CLIFlags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Load configuration
			cfg, err := config.Load(cmd)
			if err != nil {
				return err
			}

			// Setup logger
			logger := cfg.SetupLogger()

			logger.Info("Starting HackerOne Prometheus Exporter",
				slog.Int("port", int(cfg.Port)),
				slog.Int("scrape_interval", int(cfg.ScrapeInterval)),
				slog.String("log_level", cfg.LogLevel),
				slog.String("organization_id", cfg.OrganizationID),
			)

			// Create exporter
			exp := exporter.New(cfg, logger)

			// Create a new registry and register the exporter
			prometheus.MustRegister(exp)

			// Setup HTTP server
			mux := http.NewServeMux()
			mux.HandleFunc("/", handler.IndexHandler)
			mux.HandleFunc("/health", handler.HealthHandler)
			mux.Handle("/metrics", promhttp.Handler())

			server := &http.Server{
				Addr:    ":" + strconv.Itoa(int(cfg.Port)),
				Handler: mux,
			}

			// Setup context for graceful shutdown
			_, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle shutdown signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			// Start HTTP server
			go func() {
				logger.Info("HTTP server starting", slog.Int("port", int(cfg.Port)))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", slog.String("error", err.Error()))
				}
			}()

			// Wait for shutdown signal
			<-sigChan
			logger.Info("Shutdown signal received")

			// Cancel scraping context
			cancel()

			// Shutdown HTTP server
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.Error("Server shutdown error", slog.String("error", err.Error()))
			}

			logger.Info("Exporter stopped")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("Application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
