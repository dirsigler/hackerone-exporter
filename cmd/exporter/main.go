// Copyright 2025 Dennis Irsigler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
			cfg := config.New(cmd)

			// Setup logger
			logger := cfg.SetupLogger()

			logger.Info("Starting HackerOne Prometheus Exporter",
				slog.Int("port", int(cfg.Port)),
				slog.String("log_level", cfg.LogLevel),
				slog.String("organization_id", cfg.OrgID),
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
