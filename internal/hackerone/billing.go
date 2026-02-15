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

// GetProgramBalance retrieves the current bounty budget for a program.
// https://api.hackerone.com/customer-resources/#programs-get-your-programs-balance
func (c *Client) GetProgramBalance(ctx context.Context, programID string) (*types.ProgramBalance, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/billing/balance", programID)

	var balance types.ProgramBalance
	if err := c.doRequest(ctx, endpoint, &balance); err != nil {
		return nil, fmt.Errorf("getting program balance: %w", err)
	}

	return &balance, nil
}

// GetPaymentTransactions retrieves bounty payment transactions for a program.
// https://api.hackerone.com/customer-resources/#programs-get-payment-transactions
func (c *Client) GetPaymentTransactions(ctx context.Context, programID string, month, year int) (*types.PaymentTransactions, error) {
	endpoint := fmt.Sprintf("/v1/programs/%s/billing/transactions", programID)
	if month > 0 && year > 0 {
		endpoint = fmt.Sprintf("%s?month=%d&year=%d", endpoint, month, year)
	}

	var transactions types.PaymentTransactions
	if err := c.doRequest(ctx, endpoint, &transactions); err != nil {
		return nil, fmt.Errorf("getting payment transactions: %w", err)
	}

	return &transactions, nil
}
