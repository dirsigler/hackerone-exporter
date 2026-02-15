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

// GetTriageReviews retrieves triage quality reviews for a program.
// https://api.hackerone.com/customer-resources/#programs-get-triage-reviews
func (c *Client) GetTriageReviews(ctx context.Context, programID string) (*types.TriageReviews, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/triage_reviews", programID)

	var reviews types.TriageReviews
	if err := c.doRequest(ctx, endpoint, &reviews); err != nil {
		return nil, fmt.Errorf("getting triage reviews: %w", err)
	}

	return &reviews, nil
}

// GetCommonResponses retrieves template responses for a program.
// https://api.hackerone.com/customer-resources/#programs-get-common-responses
func (c *Client) GetCommonResponses(ctx context.Context, programID string) (*types.CommonResponses, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/common_responses", programID)

	var responses types.CommonResponses
	if err := c.doRequest(ctx, endpoint, &responses); err != nil {
		return nil, fmt.Errorf("getting common responses: %w", err)
	}

	return &responses, nil
}
