package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dirsigler/hackerone-exporter/pkg/types"
	"github.com/hashicorp/go-retryablehttp"
)

// HackerOneClient handles API interactions with HackerOne
type HackerOneClient struct {
	username string
	password string
	baseURL  string
	client   *retryablehttp.Client
	logger   *slog.Logger
}

// New creates a new HackerOne API client
func New(username, password, baseURL string, logger *slog.Logger) *HackerOneClient {
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = 30 * time.Second
	retryClient.Logger = nil

	return &HackerOneClient{
		username: username,
		password: password,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		client:   retryClient,
		logger:   logger,
	}
}

// makeRequest performs authenticated HTTP requests to HackerOne API
func (c *HackerOneClient) makeRequest(ctx context.Context, endpoint string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hackerone-exporter - https://github.com/dirsigler/hackerone-exporter/1.0")
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d for endpoint %s", resp.StatusCode, endpoint)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

// GetAssets retrieves all Assets for an Organization ID
// https://api.hackerone.com/customer-resources/?shell#assets-get-all-assets
func (c *HackerOneClient) GetAssets(ctx context.Context, orgID string) (*types.Assets, error) {
	var assets types.Assets
	endpoint := fmt.Sprintf("/v1/organizations/%s/assets", orgID)

	if err := c.makeRequest(ctx, endpoint, &assets); err != nil {
		return nil, fmt.Errorf("getting assets for organization %s: %w", orgID, err)
	}

	c.logger.Debug("Retrieved assets",
		slog.String("organization_id", orgID),
		slog.Int("count", len(assets.Data)))

	return &assets, nil
}

// GetAllReports retrieves all Reports for an Organization ID
// https://api.hackerone.com/customer-resources/?shell#reports-get-all-reports
func (c *HackerOneClient) GetAllReports(ctx context.Context, programHandle string) (*types.Reports, error) {
	var reports types.Reports
	endpoint := fmt.Sprintf("/v1/reports?filter[program][]=%s", programHandle)

	if err := c.makeRequest(ctx, endpoint, &reports); err != nil {
		return nil, fmt.Errorf("getting reports for program %s: %w", programHandle, err)
	}

	c.logger.Debug("Retrieved reports",
		slog.String("program", programHandle),
		slog.Int("count", len(reports.Data)))

	return &reports, nil
}

// GetPrograms retrieves all Programs
// https://api.hackerone.com/customer-resources/?shell#programs-get-your-programs
func (c *HackerOneClient) GetPrograms(ctx context.Context) (*types.Programs, error) {
	var programs types.Programs
	endpoint := "/v1/me/programs"

	if err := c.makeRequest(ctx, endpoint, &programs); err != nil {
		return nil, fmt.Errorf("getting programs: %w", err)
	}

	c.logger.Debug("Retrieved programs",
		slog.Int("count", len(programs.Data)))

	return &programs, nil
}

// GetInvitedHackers retrieves all Hackers invited to the HackerOne program
// https://api.hackerone.com/customer-resources/#programs-get-hacker-invitations
func (c *HackerOneClient) GetInvitedHackers(ctx context.Context, programID string) (*types.InvitedHackers, error) {
	var hackers types.InvitedHackers
	endpoint := fmt.Sprintf("/v1/programs/%s/hacker_invitations", programID)

	if err := c.makeRequest(ctx, endpoint, &hackers); err != nil {
		return nil, fmt.Errorf("getting hackers for program %s: %w", programID, err)
	}

	c.logger.Debug("Retrieved hackers",
		slog.String("programID", programID),
		slog.Int("count", len(hackers.Data)))

	return &hackers, nil
}

// GetWeaknesses retroves all Weaknesses for the program
// https://api.hackerone.com/customer-resources/#programs-get-weaknesses
func (c *HackerOneClient) GetWeaknesses(ctx context.Context, programID string) (*types.Weaknesses, error) {
	var weaknesses types.Weaknesses
	endpoint := fmt.Sprintf("/v1/programs/%s/weaknesses", programID)

	if err := c.makeRequest(ctx, endpoint, &weaknesses); err != nil {
		return nil, fmt.Errorf("getting weaknesses for program %s: %w", programID, err)
	}

	c.logger.Debug("Retrieved weaknesses",
		slog.String("programID", programID),
		slog.Int("count", len(weaknesses.Data)))

	return &weaknesses, nil
}
