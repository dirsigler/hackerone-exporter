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

// GetAssets retrieves all assets for an organization.
// https://api.hackerone.com/customer-resources/?shell#assets-get-all-assets
func (c *Client) GetAssets(ctx context.Context, orgID string) (*types.Assets, error) {
	endpoint := fmt.Sprintf("/v1/organizations/%s/assets", orgID)

	var assets types.Assets
	if err := c.doRequest(ctx, endpoint, &assets); err != nil {
		return nil, fmt.Errorf("getting assets for organization %s: %w", orgID, err)
	}

	return &assets, nil
}
