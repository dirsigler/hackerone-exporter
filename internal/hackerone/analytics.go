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

package hackerone

import (
	"context"
	"fmt"

	"github.com/dirsigler/hackerone-exporter/pkg/types"
)

// GetAnalytics retrieves analytics data for a given key.
// https://api.hackerone.com/customer-resources/#analytics-get-analytics-results
func (c *Client) GetAnalytics(ctx context.Context, key string) (*types.Analytics, error) {
	endpoint := fmt.Sprintf("/v1/analytics?key=%s", key)

	var analytics types.Analytics
	if err := c.doRequest(ctx, endpoint, &analytics); err != nil {
		return nil, fmt.Errorf("getting analytics for key %s: %w", key, err)
	}

	return &analytics, nil
}

// GetInvitedHackers retrieves all hackers invited to a program.
// https://api.hackerone.com/customer-resources/#programs-get-hacker-invitations
func (c *Client) GetInvitedHackers(ctx context.Context, programID string) (*types.InvitedHackers, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/hacker_invitations", programID)

	var hackers types.InvitedHackers
	if err := c.doRequest(ctx, endpoint, &hackers); err != nil {
		return nil, fmt.Errorf("getting hackers for program %s: %w", programID, err)
	}

	return &hackers, nil
}

// GetWeaknesses retrieves all weaknesses for a program.
// https://api.hackerone.com/customer-resources/#programs-get-weaknesses
func (c *Client) GetWeaknesses(ctx context.Context, programID string) (*types.Weaknesses, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/weaknesses", programID)

	var weaknesses types.Weaknesses
	if err := c.doRequest(ctx, endpoint, &weaknesses); err != nil {
		return nil, fmt.Errorf("getting weaknesses for program %s: %w", programID, err)
	}

	return &weaknesses, nil
}

// GetStructuredScopes retrieves all structured scopes for a program.
// https://api.hackerone.com/customer-resources/#programs-get-structured-scopes
func (c *Client) GetStructuredScopes(ctx context.Context, programID string) (*types.StructuredScopes, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/structured_scopes", programID)

	var scopes types.StructuredScopes
	if err := c.doRequest(ctx, endpoint, &scopes); err != nil {
		return nil, fmt.Errorf("getting structured scopes for program %s: %w", programID, err)
	}

	return &scopes, nil
}

// GetReporters retrieves all reporters for a program.
// https://api.hackerone.com/customer-resources/#programs-get-reporters
func (c *Client) GetReporters(ctx context.Context, programID string) (*types.Reporters, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/reporters", programID)

	var reporters types.Reporters
	if err := c.doRequest(ctx, endpoint, &reporters); err != nil {
		return nil, fmt.Errorf("getting reporters for program %s: %w", programID, err)
	}

	return &reporters, nil
}
