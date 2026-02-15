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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Endpoint:   "/v1/programs",
		Status:     "Not Found",
		Message:    "resource not found",
	}

	expected := "HackerOne API error: resource not found (status: 404, endpoint: /v1/programs)"
	assert.Equal(t, expected, err.Error())
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("connection refused")
	err := &NetworkError{
		URL: "/v1/programs",
		Err: originalErr,
	}

	assert.Contains(t, err.Error(), "network error for /v1/programs")
	assert.Contains(t, err.Error(), "connection refused")
	assert.Equal(t, originalErr, err.Unwrap())
}

func TestDecodeError(t *testing.T) {
	originalErr := errors.New("invalid JSON")
	err := &DecodeError{
		URL: "/v1/programs",
		Err: originalErr,
	}

	assert.Contains(t, err.Error(), "failed to decode response from /v1/programs")
	assert.Contains(t, err.Error(), "invalid JSON")
	assert.Equal(t, originalErr, err.Unwrap())
}

func TestIsAPIError(t *testing.T) {
	apiErr := &APIError{StatusCode: 404}
	assert.True(t, IsAPIError(apiErr))

	netErr := &NetworkError{URL: "/test"}
	assert.False(t, IsAPIError(netErr))

	genericErr := errors.New("generic error")
	assert.False(t, IsAPIError(genericErr))
}

func TestIsNetworkError(t *testing.T) {
	netErr := &NetworkError{URL: "/test"}
	assert.True(t, IsNetworkError(netErr))

	apiErr := &APIError{StatusCode: 404}
	assert.False(t, IsNetworkError(apiErr))

	genericErr := errors.New("generic error")
	assert.False(t, IsNetworkError(genericErr))
}

func TestIsDecodeError(t *testing.T) {
	decodeErr := &DecodeError{URL: "/test"}
	assert.True(t, IsDecodeError(decodeErr))

	apiErr := &APIError{StatusCode: 404}
	assert.False(t, IsDecodeError(apiErr))

	genericErr := errors.New("generic error")
	assert.False(t, IsDecodeError(genericErr))
}
