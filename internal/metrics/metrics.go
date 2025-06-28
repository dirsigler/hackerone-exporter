package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "hackerone"
)

// Metrics holds all Prometheus metrics for HackerOne
type Metrics struct {
	AssetsTotal         *prometheus.GaugeVec
	ReportsTotal        *prometheus.GaugeVec
	ProgramsTotal       *prometheus.GaugeVec
	InvitedHackersTotal *prometheus.GaugeVec
	WeaknessesTotal     *prometheus.GaugeVec
	LastScrapeTime      prometheus.Gauge
	ScrapeDuration      prometheus.Histogram
	ScrapeErrors        prometheus.Counter
}

var label = []string{"organization_id"}

// New creates and registers Prometheus metrics
func New() *Metrics {
	m := &Metrics{
		AssetsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "assets_total",
			Help:      "Total number of HackerOne Assets",
			Namespace: namespace,
		},
			label,
		),
		ReportsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "reports_total",
			Help:      "Total number of HackerOne Reports",
			Namespace: namespace,
		},
			label,
		),
		ProgramsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "programs_total",
			Help:      "Total number of HackerOne Programs",
			Namespace: namespace,
		},
			[]string{"handle"},
		),
		InvitedHackersTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "invited_hackers_total",
			Help:      "Total number of HackerOne Invited Hackers",
			Namespace: namespace,
		},
			append(label, "state"),
		),
		WeaknessesTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "weaknesses_total",
			Help:      "Total number of HackerOne Weaknesses",
			Namespace: namespace,
		},
			[]string{"name", "id"},
		),
		ScrapeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "scrape_errors_total",
			Help:      "Total number of HackerOne API scrape errors",
			Namespace: namespace,
		}),
		LastScrapeTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "last_scrape_timestamp",
			Help:      "Unix timestamp of the last successful scrape",
			Namespace: namespace,
		}),
		ScrapeDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:      "scrape_duration_seconds",
			Help:      "Duration of HackerOne API scrapes in seconds",
			Namespace: namespace,
			Buckets:   prometheus.DefBuckets,
		}),
	}

	return m
}

// Reset clears all metric values (useful for testing)
func (m *Metrics) Reset() {
	m.AssetsTotal.Reset()
	m.ReportsTotal.Reset()
	m.ProgramsTotal.Reset()
	m.InvitedHackersTotal.Reset()
	m.WeaknessesTotal.Reset()
	// Note: Counters and histograms cannot be reset in Prometheus
}
