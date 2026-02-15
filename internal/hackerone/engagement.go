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

// GetAllowedReporters retrieves all researchers in a private program.
// https://api.hackerone.com/customer-resources/#programs-get-allowed-reporters
func (c *Client) GetAllowedReporters(ctx context.Context, programID string) (*types.AllowedReporters, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/allowed_reporters", programID)

	var reporters types.AllowedReporters
	if err := c.doRequest(ctx, endpoint, &reporters); err != nil {
		return nil, fmt.Errorf("getting allowed reporters: %w", err)
	}

	return &reporters, nil
}

// GetAllowedReportersHistory retrieves the history of hackers joining/leaving a private program.
// https://api.hackerone.com/customer-resources/#programs-get-allowed-reporters-history
func (c *Client) GetAllowedReportersHistory(ctx context.Context, programID string) (*types.AllowedReportersHistory, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/allowed_reporters_history", programID)

	var history types.AllowedReportersHistory
	if err := c.doRequest(ctx, endpoint, &history); err != nil {
		return nil, fmt.Errorf("getting allowed reporters history: %w", err)
	}

	return &history, nil
}

// GetSwag retrieves all awarded swag for a program.
// https://api.hackerone.com/customer-resources/#programs-get-awarded-swag
func (c *Client) GetSwag(ctx context.Context, programID string) (*types.Swag, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/swag", programID)

	var swag types.Swag
	if err := c.doRequest(ctx, endpoint, &swag); err != nil {
		return nil, fmt.Errorf("getting swag: %w", err)
	}

	return &swag, nil
}
