// Copyright 2026 Dennis Irsigler
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

package collector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dirsigler/hackerone-exporter/internal/hackerone"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace         = "hackerone_reports"
	programsNamespace = "hackerone_programs"
	exporterNamespace = "hackerone_exporter"
)

// Collector implements the Prometheus Collector interface for HackerOne metrics.
type Collector struct {
	// Report metrics
	totalCount          *prometheus.Desc
	stateCount          *prometheus.Desc
	severityCount       *prometheus.Desc
	reporterReportCount *prometheus.Desc
	timeToTriage        *prometheus.Desc
	timeToResolution    *prometheus.Desc

	// Program metrics
	programsTotal *prometheus.Desc

	// Financial metrics
	programBalance          *prometheus.Desc
	transactionsTotal       *prometheus.Desc
	transactionsBountyTotal *prometheus.Desc

	// Engagement metrics
	allowedReportersTotal         *prometheus.Desc
	allowedReportersVerifiedTotal *prometheus.Desc
	swagTotal                     *prometheus.Desc
	swagSent                      *prometheus.Desc
	swagPending                   *prometheus.Desc

	// Quality metrics
	triageReviewsTotal    *prometheus.Desc
	triageReviewRatingAvg *prometheus.Desc
	commonResponsesTotal  *prometheus.Desc

	// Health metrics
	reportersTotal  *prometheus.Desc
	thanksTotal     *prometheus.Desc
	weaknessesTotal *prometheus.Desc

	// Team metrics
	programMembersTotal *prometheus.Desc
	programAdminsTotal  *prometheus.Desc

	// Exporter metrics
	up                    prometheus.Gauge
	collectionDuration    prometheus.Gauge
	apiCallDuration       *prometheus.GaugeVec
	apiCallErrors         *prometheus.CounterVec
	lastCollectionSuccess prometheus.Gauge
	lastCollectionAttempt prometheus.Gauge
	collectionErrors      prometheus.Counter

	client         *hackerone.Client
	programHandles []string
	mu             sync.RWMutex
	lastSuccess    time.Time
	lastError      error
}

// New creates a new HackerOne Prometheus collector.
func New(client *hackerone.Client, programHandles []string) *Collector {
	return &Collector{
		// Report metrics
		totalCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "total", "count"),
			"The total number of reports.",
			nil,
			nil,
		),
		stateCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "state", "count"),
			"The number of reports by state.",
			[]string{"state"},
			nil,
		),
		severityCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "severity", "count"),
			"The number of reports by severity.",
			[]string{"severity"},
			nil,
		),
		reporterReportCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "reporter", "count"),
			"The number of reports submitted by each reporter.",
			[]string{"reporter_username", "reporter_id"},
			nil,
		),
		timeToTriage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "time_to_triage", "seconds"),
			"Time from report creation to triage in seconds.",
			nil,
			nil,
		),
		timeToResolution: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "time_to_resolution", "seconds"),
			"Time from report creation to close/resolution in seconds.",
			nil,
			nil,
		),
		programsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "total", "count"),
			"The total number of programs.",
			nil,
			nil,
		),
		programBalance: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "balance", "usd"),
			"Current bounty balance for the program in USD.",
			[]string{"program_id"},
			nil,
		),
		transactionsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "transactions", "total"),
			"Total number of payment transactions.",
			[]string{"program_id"},
			nil,
		),
		transactionsBountyTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "transactions_bounty", "usd_total"),
			"Total bounty amount paid in USD.",
			[]string{"program_id"},
			nil,
		),
		// Engagement metrics
		allowedReportersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "allowed_reporters", "total"),
			"Total number of allowed reporters in private program.",
			[]string{"program_id"},
			nil,
		),
		allowedReportersVerifiedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "allowed_reporters_verified", "total"),
			"Number of identity-verified allowed reporters.",
			[]string{"program_id"},
			nil,
		),
		swagTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "swag", "total"),
			"Total number of swag items awarded.",
			[]string{"program_id"},
			nil,
		),
		swagSent: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "swag_sent", "total"),
			"Number of swag items marked as sent.",
			[]string{"program_id"},
			nil,
		),
		swagPending: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "swag_pending", "total"),
			"Number of swag items pending shipment.",
			[]string{"program_id"},
			nil,
		),
		// Quality metrics
		triageReviewsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "triage_reviews", "total"),
			"Total number of triage reviews received.",
			[]string{"program_id"},
			nil,
		),
		triageReviewRatingAvg: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "triage_review_rating", "average"),
			"Average triage review rating (1-5 scale).",
			[]string{"program_id"},
			nil,
		),
		commonResponsesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "common_responses", "total"),
			"Total number of common response templates.",
			[]string{"program_id"},
			nil,
		),
		// Health metrics
		reportersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "reporters", "total"),
			"Total number of unique reporters who submitted.",
			[]string{"program_id"},
			nil,
		),
		thanksTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "thanks", "total"),
			"Total number of hackers in thanks list.",
			[]string{"program_id"},
			nil,
		),
		weaknessesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "weaknesses", "total"),
			"Total number of weakness categories available.",
			[]string{"program_id"},
			nil,
		),
		// Team metrics
		programMembersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "members", "total"),
			"Total number of team members with program access.",
			[]string{"program_id"},
			nil,
		),
		programAdminsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(programsNamespace, "admins", "total"),
			"Total number of organization admins.",
			[]string{"program_id"},
			nil,
		),
		// Exporter metrics
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: exporterNamespace,
			Name:      "up",
			Help:      "Was the last scrape of the HackerOne API successful.",
		}),
		collectionDuration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: exporterNamespace,
			Name:      "collection_duration_seconds",
			Help:      "Duration of the last collection cycle in seconds.",
		}),
		apiCallDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: exporterNamespace,
			Name:      "api_call_duration_seconds",
			Help:      "Duration of API calls in seconds.",
		}, []string{"endpoint"}),
		apiCallErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: exporterNamespace,
			Name:      "api_call_errors_total",
			Help:      "Total number of API call errors.",
		}, []string{"endpoint"}),
		lastCollectionSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: exporterNamespace,
			Name:      "last_collection_success_timestamp_seconds",
			Help:      "Timestamp of the last successful collection.",
		}),
		lastCollectionAttempt: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: exporterNamespace,
			Name:      "last_collection_attempt_timestamp_seconds",
			Help:      "Timestamp of the last collection attempt.",
		}),
		collectionErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: exporterNamespace,
			Name:      "collection_errors_total",
			Help:      "Total number of collection errors.",
		}),
		client:         client,
		programHandles: programHandles,
	}
}

// Describe sends the super-set of all possible descriptors to Prometheus.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Report metrics
	ch <- c.totalCount
	ch <- c.stateCount
	ch <- c.severityCount
	ch <- c.reporterReportCount
	ch <- c.timeToTriage
	ch <- c.timeToResolution

	// Program metrics
	ch <- c.programsTotal

	// Financial metrics
	ch <- c.programBalance
	ch <- c.transactionsTotal
	ch <- c.transactionsBountyTotal

	// Engagement metrics
	ch <- c.allowedReportersTotal
	ch <- c.allowedReportersVerifiedTotal
	ch <- c.swagTotal
	ch <- c.swagSent
	ch <- c.swagPending

	// Quality metrics
	ch <- c.triageReviewsTotal
	ch <- c.triageReviewRatingAvg
	ch <- c.commonResponsesTotal

	// Health metrics
	ch <- c.reportersTotal
	ch <- c.thanksTotal
	ch <- c.weaknessesTotal

	// Team metrics
	ch <- c.programMembersTotal
	ch <- c.programAdminsTotal

	// Exporter metrics
	c.up.Describe(ch)
	c.collectionDuration.Describe(ch)
	c.apiCallDuration.Describe(ch)
	c.apiCallErrors.Describe(ch)
	c.lastCollectionSuccess.Describe(ch)
	c.lastCollectionAttempt.Describe(ch)
	c.collectionErrors.Describe(ch)
}

// Collect fetches metrics from HackerOne and sends them to Prometheus.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	slog.Info("starting to collect metrics")
	start := time.Now()
	c.lastCollectionAttempt.SetToCurrentTime()

	// Create context with timeout for all API calls
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Fetch all programs
	apiStart := time.Now()
	programs, err := c.client.GetPrograms(ctx)
	c.apiCallDuration.WithLabelValues("programs").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Error("failed to get programs, aborting collection",
			"error", err,
			"component", "collector",
		)
		c.apiCallErrors.WithLabelValues("programs").Inc()
		c.collectionErrors.Inc()
		c.up.Set(0)
		c.mu.Lock()
		c.lastError = err
		c.mu.Unlock()

		// Collect exporter metrics even on failure
		c.up.Collect(ch)
		c.collectionDuration.Collect(ch)
		c.apiCallDuration.Collect(ch)
		c.apiCallErrors.Collect(ch)
		c.lastCollectionSuccess.Collect(ch)
		c.lastCollectionAttempt.Collect(ch)
		c.collectionErrors.Collect(ch)
		return
	}

	slog.Debug("programs retrieved", "count", len(programs.Data))

	// Emit total programs count
	ch <- prometheus.MustNewConstMetric(
		c.programsTotal,
		prometheus.GaugeValue,
		float64(len(programs.Data)),
	)

	// Collect reports for configured programs (or all if none specified)
	programsToScrape := c.programHandles
	if len(programsToScrape) == 0 {
		// If no specific programs configured, scrape all
		for _, prog := range programs.Data {
			programsToScrape = append(programsToScrape, prog.Attributes.Handle)
		}
	}

	// Histogram buckets: 1h, 4h, 8h, 24h, 3d, 7d, 14d, 30d, 60d, 90d
	durationBuckets := []float64{
		3600, 14400, 28800, 86400,
		259200, 604800, 1209600,
		2592000, 5184000, 7776000,
	}

	// Aggregate reports from all programs
	allReports := make([]interface{}, 0)
	stateCounts := make(map[string]int)
	var triageDurations []float64
	var resolutionDurations []float64

	// Track reports by reporter
	type reporterInfo struct {
		username string
		count    int
	}
	reporterCounts := make(map[string]*reporterInfo) // map[reporter_id]*reporterInfo

	for _, programHandle := range programsToScrape {
		apiStart := time.Now()
		reports, err := c.client.GetAllReports(ctx, programHandle)
		c.apiCallDuration.WithLabelValues("reports").Set(time.Since(apiStart).Seconds())

		if err != nil {
			slog.Error("failed to get reports for program",
				"error", err,
				"program", programHandle,
				"component", "collector",
			)
			c.apiCallErrors.WithLabelValues("reports").Inc()
			continue
		}

		slog.Debug("reports retrieved",
			"program", programHandle,
			"count", len(reports.Data),
		)

		// Count reports by state and reporter
		for _, report := range reports.Data {
			allReports = append(allReports, report)

			// Count by state
			if report.Attributes.State != "" {
				stateCounts[report.Attributes.State]++
			}

			// Count by reporter
			if report.Relationships.Reporter.Data.ID != "" {
				reporterID := report.Relationships.Reporter.Data.ID
				reporterUsername := report.Relationships.Reporter.Data.Attributes.Username

				if _, exists := reporterCounts[reporterID]; !exists {
					reporterCounts[reporterID] = &reporterInfo{
						username: reporterUsername,
						count:    0,
					}
				}
				reporterCounts[reporterID].count++
			}

			// Collect time-to-triage durations
			if report.Attributes.TriagedAt != nil {
				triageDurations = append(triageDurations, report.Attributes.TriagedAt.Sub(report.Attributes.CreatedAt).Seconds())
			}

			// Collect time-to-resolution durations
			if report.Attributes.ClosedAt != nil {
				resolutionDurations = append(resolutionDurations, report.Attributes.ClosedAt.Sub(report.Attributes.CreatedAt).Seconds())
			}
		}
	}

	// Emit total count metric
	ch <- prometheus.MustNewConstMetric(
		c.totalCount,
		prometheus.GaugeValue,
		float64(len(allReports)),
	)

	// Emit state metrics
	// Common HackerOne report states
	commonStates := []string{"new", "triaged", "needs-more-info", "resolved", "not-applicable", "informative", "duplicate", "spam"}
	for _, state := range commonStates {
		count := stateCounts[state]
		slog.Debug("emitting state metric",
			"state", state,
			"count", count,
		)
		ch <- prometheus.MustNewConstMetric(
			c.stateCount,
			prometheus.GaugeValue,
			float64(count),
			state,
		)
	}

	// Emit reporter metrics
	for reporterID, info := range reporterCounts {
		slog.Debug("emitting reporter metric",
			"reporter_id", reporterID,
			"reporter_username", info.username,
			"count", info.count,
		)
		ch <- prometheus.MustNewConstMetric(
			c.reporterReportCount,
			prometheus.GaugeValue,
			float64(info.count),
			info.username,
			reporterID,
		)
	}

	// Emit time-to-triage histogram
	ch <- buildConstHistogram(c.timeToTriage, durationBuckets, triageDurations)

	// Emit time-to-resolution histogram
	ch <- buildConstHistogram(c.timeToResolution, durationBuckets, resolutionDurations)

	// Collect program-specific metrics for each program
	for _, prog := range programs.Data {
		programID := prog.ID
		slog.Debug("collecting program-specific metrics", "program_id", programID)

		// Financial metrics
		c.collectFinancialMetrics(ctx, programID, ch)

		// Engagement metrics
		c.collectEngagementMetrics(ctx, programID, ch)

		// Quality metrics
		c.collectQualityMetrics(ctx, programID, ch)

		// Health metrics
		c.collectHealthMetrics(ctx, programID, ch)

		// Team metrics
		c.collectTeamMetrics(ctx, programID, ch)
	}

	duration := time.Since(start)
	c.collectionDuration.Set(duration.Seconds())
	c.lastCollectionSuccess.SetToCurrentTime()
	c.up.Set(1)

	c.mu.Lock()
	c.lastSuccess = time.Now()
	c.lastError = nil
	c.mu.Unlock()

	// Collect exporter metrics
	c.up.Collect(ch)
	c.collectionDuration.Collect(ch)
	c.apiCallDuration.Collect(ch)
	c.apiCallErrors.Collect(ch)
	c.lastCollectionSuccess.Collect(ch)
	c.lastCollectionAttempt.Collect(ch)
	c.collectionErrors.Collect(ch)

	slog.Info("finished collecting metrics",
		"duration_seconds", duration.Seconds(),
		"component", "collector",
	)
}

// HealthStatus returns the health status of the collector.
func (c *Collector) HealthStatus() (healthy bool, lastSuccess time.Time, lastError error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Consider healthy if we've had a successful collection or no error yet
	healthy = c.lastError == nil
	return healthy, c.lastSuccess, c.lastError
}

// CheckUpstreamHealth performs a lightweight health check against the HackerOne API.
func (c *Collector) CheckUpstreamHealth() error {
	return c.client.HealthCheck()
}

// collectFinancialMetrics collects balance and transaction metrics for a program.
func (c *Collector) collectFinancialMetrics(ctx context.Context, programID string, ch chan<- prometheus.Metric) {
	// Get program balance
	apiStart := time.Now()
	balance, err := c.client.GetProgramBalance(ctx, programID)
	c.apiCallDuration.WithLabelValues("program_balance").Set(time.Since(apiStart).Seconds())

	if err != nil {
		// Check if this is a permission or configuration error
		var apiErr *hackerone.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 403:
				slog.Debug("skipping program balance (insufficient permissions)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			case 401:
				slog.Debug("skipping program balance (feature not configured/available)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			}
		}
		slog.Warn("failed to get program balance",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("program_balance").Inc()
		return
	}

	// Parse balance as float (format is "1500.00")
	var balanceValue float64
	if _, parseErr := fmt.Sscanf(balance.Data.Attributes.Balance, "%f", &balanceValue); parseErr == nil {
		ch <- prometheus.MustNewConstMetric(
			c.programBalance,
			prometheus.GaugeValue,
			balanceValue,
			programID,
		)
	}

	// Get payment transactions (current month)
	now := time.Now()
	apiStart = time.Now()
	transactions, err := c.client.GetPaymentTransactions(ctx, programID, int(now.Month()), now.Year())
	c.apiCallDuration.WithLabelValues("payment_transactions").Set(time.Since(apiStart).Seconds())

	if err != nil {
		// Check if this is a permission or configuration error
		var apiErr *hackerone.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 403:
				slog.Debug("skipping payment transactions (insufficient permissions)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			case 401:
				slog.Debug("skipping payment transactions (feature not configured/available)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			}
		}
		slog.Warn("failed to get payment transactions",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("payment_transactions").Inc()
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.transactionsTotal,
		prometheus.GaugeValue,
		float64(len(transactions.Data)),
		programID,
	)

	// Sum up total bounty amounts
	var totalBounty float64
	for _, tx := range transactions.Data {
		var bountyAmount float64
		if _, parseErr := fmt.Sscanf(tx.BountyAward, "%f", &bountyAmount); parseErr == nil {
			totalBounty += bountyAmount
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.transactionsBountyTotal,
		prometheus.GaugeValue,
		totalBounty,
		programID,
	)
}

// collectEngagementMetrics collects allowed reporters and swag metrics for a program.
func (c *Collector) collectEngagementMetrics(ctx context.Context, programID string, ch chan<- prometheus.Metric) {
	// Get allowed reporters
	apiStart := time.Now()
	reporters, err := c.client.GetAllowedReporters(ctx, programID)
	c.apiCallDuration.WithLabelValues("allowed_reporters").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Debug("failed to get allowed reporters (may not be a private program)",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("allowed_reporters").Inc()
	} else {
		verifiedCount := 0
		for _, reporter := range reporters.Data {
			if reporter.Attributes.IdentityVerified {
				verifiedCount++
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.allowedReportersTotal,
			prometheus.GaugeValue,
			float64(len(reporters.Data)),
			programID,
		)

		ch <- prometheus.MustNewConstMetric(
			c.allowedReportersVerifiedTotal,
			prometheus.GaugeValue,
			float64(verifiedCount),
			programID,
		)
	}

	// Get swag
	apiStart = time.Now()
	swag, err := c.client.GetSwag(ctx, programID)
	c.apiCallDuration.WithLabelValues("swag").Set(time.Since(apiStart).Seconds())

	if err != nil {
		// Check if this is a permission or configuration error
		var apiErr *hackerone.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 403:
				slog.Debug("skipping swag metrics (insufficient permissions)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			case 401:
				slog.Debug("skipping swag metrics (feature not configured/available)",
					"status_code", apiErr.StatusCode,
					"program_id", programID,
					"error", err,
				)
				return
			}
		}
		slog.Warn("failed to get swag",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("swag").Inc()
		return
	}

	sentCount := 0
	pendingCount := 0
	for _, item := range swag.Data {
		if item.Attributes.Sent {
			sentCount++
		} else {
			pendingCount++
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.swagTotal,
		prometheus.GaugeValue,
		float64(len(swag.Data)),
		programID,
	)

	ch <- prometheus.MustNewConstMetric(
		c.swagSent,
		prometheus.GaugeValue,
		float64(sentCount),
		programID,
	)

	ch <- prometheus.MustNewConstMetric(
		c.swagPending,
		prometheus.GaugeValue,
		float64(pendingCount),
		programID,
	)
}

// collectQualityMetrics collects triage reviews and common responses metrics for a program.
func (c *Collector) collectQualityMetrics(ctx context.Context, programID string, ch chan<- prometheus.Metric) {
	// Get triage reviews
	apiStart := time.Now()
	reviews, err := c.client.GetTriageReviews(ctx, programID)
	c.apiCallDuration.WithLabelValues("triage_reviews").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get triage reviews",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("triage_reviews").Inc()
	} else {
		totalRating := 0
		for _, review := range reviews.Data {
			totalRating += review.Attributes.Rating
		}

		ch <- prometheus.MustNewConstMetric(
			c.triageReviewsTotal,
			prometheus.GaugeValue,
			float64(len(reviews.Data)),
			programID,
		)

		if len(reviews.Data) > 0 {
			avgRating := float64(totalRating) / float64(len(reviews.Data))
			ch <- prometheus.MustNewConstMetric(
				c.triageReviewRatingAvg,
				prometheus.GaugeValue,
				avgRating,
				programID,
			)
		}
	}

	// Get common responses
	apiStart = time.Now()
	responses, err := c.client.GetCommonResponses(ctx, programID)
	c.apiCallDuration.WithLabelValues("common_responses").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get common responses",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("common_responses").Inc()
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.commonResponsesTotal,
		prometheus.GaugeValue,
		float64(len(responses.Data)),
		programID,
	)
}

// collectHealthMetrics collects reporters, thanks, and weaknesses metrics for a program.
func (c *Collector) collectHealthMetrics(ctx context.Context, programID string, ch chan<- prometheus.Metric) {
	// Get reporters
	apiStart := time.Now()
	reporters, err := c.client.GetReporters(ctx, programID)
	c.apiCallDuration.WithLabelValues("reporters").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get reporters",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("reporters").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.reportersTotal,
			prometheus.GaugeValue,
			float64(len(reporters.Data)),
			programID,
		)
	}

	// Get thanks
	apiStart = time.Now()
	thanks, err := c.client.GetThanks(ctx, programID)
	c.apiCallDuration.WithLabelValues("thanks").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get thanks",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("thanks").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.thanksTotal,
			prometheus.GaugeValue,
			float64(len(thanks.Data)),
			programID,
		)
	}

	// Get weaknesses
	apiStart = time.Now()
	weaknesses, err := c.client.GetWeaknesses(ctx, programID)
	c.apiCallDuration.WithLabelValues("weaknesses").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get weaknesses",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("weaknesses").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.weaknessesTotal,
			prometheus.GaugeValue,
			float64(len(weaknesses.Data)),
			programID,
		)
	}
}

// collectTeamMetrics collects program members metrics for a program.
func (c *Collector) collectTeamMetrics(ctx context.Context, programID string, ch chan<- prometheus.Metric) {
	// Get program members
	apiStart := time.Now()
	members, err := c.client.GetProgramMembers(ctx, programID)
	c.apiCallDuration.WithLabelValues("program_members").Set(time.Since(apiStart).Seconds())

	if err != nil {
		slog.Warn("failed to get program members",
			"error", err,
			"program_id", programID,
		)
		c.apiCallErrors.WithLabelValues("program_members").Inc()
		return
	}

	adminCount := 0
	for _, member := range members.Data {
		if member.Attributes.OrganizationAdmin {
			adminCount++
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.programMembersTotal,
		prometheus.GaugeValue,
		float64(len(members.Data)),
		programID,
	)

	ch <- prometheus.MustNewConstMetric(
		c.programAdminsTotal,
		prometheus.GaugeValue,
		float64(adminCount),
		programID,
	)
}

// buildConstHistogram creates a constant histogram metric from observed values.
func buildConstHistogram(desc *prometheus.Desc, buckets []float64, observations []float64) prometheus.Metric {
	bucketCounts := make(map[float64]uint64, len(buckets))
	for _, b := range buckets {
		bucketCounts[b] = 0
	}

	var sum float64
	for _, v := range observations {
		sum += v
		for _, b := range buckets {
			if v <= b {
				bucketCounts[b]++
			}
		}
	}

	return prometheus.MustNewConstHistogram(
		desc,
		uint64(len(observations)),
		sum,
		bucketCounts,
	)
}
