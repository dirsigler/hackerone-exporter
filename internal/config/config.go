package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/urfave/cli/v3"
)

// Config holds the application configuration
type Config struct {
	HackerOneBasicAuthUsername string `env:"HACKERONE_BASIC_AUTH_USERNAME,required"`
	HackerOneBasicAuthPassword string `env:"HACKERONE_BASIC_AUTH_PASSWORD,required"`
	Port                       int64  `env:"PORT" envDefault:"8080"`
	ScrapeInterval             int64  `env:"SCRAPE_INTERVAL" envDefault:"60"`
	LogLevel                   string `env:"LOG_LEVEL" envDefault:"info"`
	HackerOneAPIURL            string `env:"HACKERONE_API_URL" envDefault:"https://api.hackerone.com"`
	OrganizationID             string `env:"HACKERONE_ORGANIZATION_ID,required"`
}

// Load parses configuration from environment variables and CLI context
func Load(cmd *cli.Command) (*Config, error) {
	// Parse configuration from environment variables
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("parsing config from environment: %w", err)
	}

	// Override with CLI flags if provided
	if cmd.String("hackerone-basic-auth-username") != "" {
		config.HackerOneBasicAuthUsername = cmd.String("hackerone-basic-auth-username")
	}
	if cmd.String("hackerone-basic-auth-password") != "" {
		config.HackerOneBasicAuthPassword = cmd.String("hackerone-basic-auth-password")
	}
	if cmd.Int("port") != 8080 {
		config.Port = cmd.Int("port")
	}
	if cmd.Int("scrape-interval") != 60 {
		config.ScrapeInterval = cmd.Int("scrape-interval")
	}
	if cmd.String("log-level") != "info" {
		config.LogLevel = cmd.String("log-level")
	}
	if cmd.String("organization-id") != "" {
		config.OrganizationID = cmd.String("organization-id")
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.HackerOneBasicAuthUsername == "" {
		return fmt.Errorf("HackerOne Basic Auth Username is required")
	}
	if c.HackerOneBasicAuthPassword == "" {
		return fmt.Errorf("HackerOne Basic Auth Password is required")
	}
	if c.OrganizationID == "" {
		return fmt.Errorf("HackerOne Organization ID is required")
	}

	return nil
}

// SetupLogger configures the structured logger based on config
func (c *Config) SetupLogger() *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// CLIFlags returns the CLI flags for the application
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "hackerone-basic-auth-username",
			Usage:   "HackerOne Basic Auth Username (can also be set via HACCERONE_BASIC_AUTH_USERNAME env var)",
			Sources: cli.EnvVars("HACKERONE_BASIC_AUTH_USERNAME"),
		},
		&cli.StringFlag{
			Name:    "hackerone-basic-auth-password",
			Usage:   "HackerOne Basic Auth Password (can also be set via HACCERONE_BASIC_AUTH_PASSWORD env var)",
			Sources: cli.EnvVars("HACKERONE_BASIC_AUTH_PASSWORD"),
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port to listen on",
			Sources: cli.EnvVars("PORT"),
			Value:   8080,
		},
		&cli.IntFlag{
			Name:    "scrape-interval",
			Usage:   "Scrape interval in seconds",
			Sources: cli.EnvVars("SCRAPE_INTERVAL"),
			Value:   60,
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level (debug, info, warn, error)",
			Sources: cli.EnvVars("LOG_LEVEL"),
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "organization-id",
			Usage:   "HackerOne Organization ID",
			Sources: cli.EnvVars("HACKERONE_ORGANIZATION_ID"),
		},
	}
}
