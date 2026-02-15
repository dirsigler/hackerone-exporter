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

package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Config holds the entire application configuration.
type Config struct {
	HackerOne HackerOneConfig
	Port      int
	Level     slog.Level
}

// HackerOneConfig holds HackerOne API configuration.
type HackerOneConfig struct {
	Username string
	Password string
	URL      string
	OrgID string
}

// LoadParams holds parameters for configuration loading.
type LoadParams struct {
	Port        int
	LogLevel    string
	APIUser     string
	APIPassword string
	APIURL      string
	OrgID       string
}

// Load builds the configuration from provided parameters.
func Load(params LoadParams) (*Config, error) {
	var cfg Config

	// Parse log level
	var level slog.Level
	switch strings.ToLower(params.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", params.LogLevel)
	}
	cfg.Level = level
	cfg.Port = params.Port

	// Validate required fields
	if params.APIUser == "" {
		return nil, errors.New("API user is required")
	}
	cfg.HackerOne.Username = params.APIUser

	if params.APIPassword == "" {
		return nil, errors.New("API password is required")
	}
	cfg.HackerOne.Password = params.APIPassword

	if params.OrgID == "" {
		return nil, errors.New("organization ID is required")
	}
	cfg.HackerOne.OrgID = params.OrgID

	// Set API URL
	if params.APIURL == "" {
		params.APIURL = "https://api.hackerone.com"
	}
	cfg.HackerOne.URL = params.APIURL

	return &cfg, nil
}
