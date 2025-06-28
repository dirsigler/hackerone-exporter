// Copyright 2025 Dennis Irsigler
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
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

// Config holds the application configuration
type Config struct {
	APIUser     string
	APIPassword string
	Port        int64
	LogLevel    string
	APIURL      string
	OrgID       string
}

// New creates a new Config struct from the cli.Command
func New(cmd *cli.Command) *Config {
	return &Config{
		APIUser:     cmd.String("api-user"),
		APIPassword: cmd.String("api-password"),
		Port:        cmd.Int("port"),
		LogLevel:    cmd.String("log-level"),
		APIURL:      cmd.String("api-url"),
		OrgID:       cmd.String("org-id"),
	}
}

// SetupLogger configures the structured logger based on config
func (c *Config) SetupLogger() *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(c.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// CLIFlags returns the CLI flags for the application
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "api-user",
			Usage:    "HackerOne API Username",
			Sources:  cli.EnvVars("HACKERONE_API_USER"),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "api-password",
			Usage:    "HackerOne API Password",
			Sources:  cli.EnvVars("HACKERONE_API_PASSWORD"),
			Required: true,
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port to listen on",
			Sources: cli.EnvVars("PORT"),
			Value:   8080,
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level (debug, info, warn, error)",
			Sources: cli.EnvVars("LOG_LEVEL"),
			Value:   "info",
		},
		&cli.StringFlag{
			Name:     "org-id",
			Usage:    "HackerOne Organization ID",
			Sources:  cli.EnvVars("HACKERONE_ORG_ID"),
			Required: true,
		},
		&cli.StringFlag{
			Name:    "api-url",
			Usage:   "HackerOne API URL",
			Sources: cli.EnvVars("HACKERONE_API_URL"),
			Value:   "https://api.hackerone.com",
			Hidden:  true,
		},
	}
}
