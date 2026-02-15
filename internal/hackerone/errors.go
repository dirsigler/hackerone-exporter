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
	"errors"
	"fmt"
)

// APIError represents an error from the HackerOne API.
type APIError struct {
	StatusCode int
	Endpoint   string
	Status     string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HackerOne API error: %s (status: %d, endpoint: %s)",
		e.Message, e.StatusCode, e.Endpoint)
}

// NetworkError represents a network-related error during API calls.
type NetworkError struct {
	URL string
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error for %s: %v", e.URL, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// DecodeError represents an error decoding the API response.
type DecodeError struct {
	URL string
	Err error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("failed to decode response from %s: %v", e.URL, e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

// IsAPIError checks if an error is an APIError.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsNetworkError checks if an error is a NetworkError.
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// IsDecodeError checks if an error is a DecodeError.
func IsDecodeError(err error) bool {
	var decodeErr *DecodeError
	return errors.As(err, &decodeErr)
}
