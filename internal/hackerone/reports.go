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

// GetAllReports retrieves all reports for a program.
// https://api.hackerone.com/customer-resources/?shell#reports-get-all-reports
func (c *Client) GetAllReports(ctx context.Context, programHandle string) (*types.Reports, error) {
	endpoint := fmt.Sprintf("/v1/reports?filter[program][]=%s", programHandle)

	var reports types.Reports
	if err := c.doRequest(ctx, endpoint, &reports); err != nil {
		return nil, fmt.Errorf("getting reports for program %s: %w", programHandle, err)
	}

	return &reports, nil
}
