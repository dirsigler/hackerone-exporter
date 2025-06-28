package exporter

import (
	"context"

	"log/slog"
	"sync"
	"time"

	"github.com/dirsigler/hackerone-exporter/internal/client"
	"github.com/dirsigler/hackerone-exporter/internal/config"
	"github.com/dirsigler/hackerone-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter manages the HackerOne metrics collection
type Exporter struct {
	client  *client.HackerOneClient
	metrics *metrics.Metrics
	logger  *slog.Logger
	config  *config.Config
	mu      sync.RWMutex
}

// New creates a new HackerOne exporter
func New(cfg *config.Config, logger *slog.Logger) *Exporter {
	HackerOneClient := client.New(cfg.HackerOneBasicAuthUsername, cfg.HackerOneBasicAuthPassword, cfg.HackerOneAPIURL, logger)
	prometheusMetrics := metrics.New()

	return &Exporter{
		client:  HackerOneClient,
		metrics: prometheusMetrics,
		logger:  logger,
		config:  cfg,
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// that can be collected by this Collector to the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.metrics.AssetsTotal.Describe(ch)
	e.metrics.ReportsTotal.Describe(ch)
	e.metrics.ProgramsTotal.Describe(ch)
	e.metrics.InvitedHackersTotal.Describe(ch)
	e.metrics.WeaknessesTotal.Describe(ch)
	e.metrics.ScrapeErrors.Describe(ch)
	e.metrics.LastScrapeTime.Describe(ch)
	e.metrics.ScrapeDuration.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.metrics.Reset()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.config.ScrapeInterval)*time.Second)
	defer cancel()

	timer := prometheus.NewTimer(e.metrics.ScrapeDuration)
	defer timer.ObserveDuration()

	e.logger.Info("Starting HackerOne metrics scrape")

	assets, err := e.client.GetAssets(ctx, e.config.OrganizationID)
	if err != nil {
		e.metrics.ScrapeErrors.Inc()
		e.logger.Error("getting assets", slog.String("error", err.Error()))
	}

	programs, err := e.client.GetPrograms(ctx)
	if err != nil {
		e.metrics.ScrapeErrors.Inc()
		e.logger.Error("getting programs", slog.String("error", err.Error()))
	}

	for _, program := range programs.Data {
		e.metrics.ProgramsTotal.WithLabelValues(program.Attributes.Handle).Inc()

		reports, err := e.client.GetAllReports(ctx, program.Attributes.Handle)
		if err != nil {
			e.metrics.ScrapeErrors.Inc()
			e.logger.Error("getting reports for program", slog.String("program", program.ID), slog.String("error", err.Error()))
		}
		e.metrics.ReportsTotal.WithLabelValues(e.config.OrganizationID).Set(float64(len(reports.Data)))

		hackers, err := e.client.GetInvitedHackers(ctx, program.ID)
		if err != nil {
			e.metrics.ScrapeErrors.Inc()
			e.logger.Error("getting hackers for program", slog.String("program", program.ID), slog.String("error", err.Error()))
		}
		for _, hacker := range hackers.Data {
			e.metrics.InvitedHackersTotal.WithLabelValues(e.config.OrganizationID, hacker.Attributes.State).Inc()
		}

		weaknesses, err := e.client.GetWeaknesses(ctx, program.ID)
		if err != nil {
			e.metrics.ScrapeErrors.Inc()
			e.logger.Error("getting weaknesses for program", slog.String("program", program.ID), slog.String("error", err.Error()))
		}
		for _, weakness := range weaknesses.Data {
			e.metrics.WeaknessesTotal.WithLabelValues(weakness.Attributes.Name, weakness.ID).Inc()
		}

	}
	e.metrics.AssetsTotal.WithLabelValues(e.config.OrganizationID).Set(float64(len(assets.Data)))

	e.metrics.LastScrapeTime.SetToCurrentTime()
	e.logger.Info("HackerOne metrics scrape completed")

	e.metrics.AssetsTotal.Collect(ch)
	e.metrics.ReportsTotal.Collect(ch)
	e.metrics.ProgramsTotal.Collect(ch)
	e.metrics.InvitedHackersTotal.Collect(ch)
	e.metrics.WeaknessesTotal.Collect(ch)
	e.metrics.ScrapeErrors.Collect(ch)
	e.metrics.LastScrapeTime.Collect(ch)
	e.metrics.ScrapeDuration.Collect(ch)
}
