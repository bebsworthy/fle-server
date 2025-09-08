// Package config provides configuration management for the FLE server.
// It handles loading configuration from environment variables with sensible defaults,
// supporting both development and production environments.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Default configuration constants
const (
	DefaultPort                     = 8080
	DefaultHost                     = "0.0.0.0"
	DefaultCORSOrigin               = "http://localhost:3000"
	DefaultLogLevel                 = "info"
	DefaultEnvironment              = "development"
	DefaultWebSocketReadBufferSize  = 1024
	DefaultWebSocketWriteBufferSize = 1024
	DefaultMaxConnections           = 1000
	DefaultHeartbeatInterval        = 30   // seconds
	DefaultSessionTimeout           = 3600 // 1 hour in seconds
)

// Config represents the complete configuration for the FLE server.
// All fields can be configured via environment variables with fallback defaults.
type Config struct {
	// Server configuration
	Port int    `json:"port" env:"PORT"`
	Host string `json:"host" env:"HOST"`

	// CORS configuration for frontend development
	CORSOrigin string `json:"corsOrigin" env:"CORS_ORIGIN"`

	// Logging configuration
	LogLevel string `json:"logLevel" env:"LOG_LEVEL"`

	// Environment (development, production, test)
	Environment string `json:"environment" env:"ENV"`

	// WebSocket configuration
	WebSocketReadBufferSize  int `json:"wsReadBufferSize" env:"WS_READ_BUFFER_SIZE"`
	WebSocketWriteBufferSize int `json:"wsWriteBufferSize" env:"WS_WRITE_BUFFER_SIZE"`

	// Connection management
	MaxConnections    int `json:"maxConnections" env:"MAX_CONNECTIONS"`
	HeartbeatInterval int `json:"heartbeatInterval" env:"HEARTBEAT_INTERVAL"`

	// Session configuration
	SessionTimeout int `json:"sessionTimeout" env:"SESSION_TIMEOUT"`
}

// defaultConfig returns the default configuration values.
// These are production-safe defaults that work well for development.
func defaultConfig() *Config {
	return &Config{
		Port:                     DefaultPort,
		Host:                     DefaultHost,
		CORSOrigin:               DefaultCORSOrigin,
		LogLevel:                 DefaultLogLevel,
		Environment:              DefaultEnvironment,
		WebSocketReadBufferSize:  DefaultWebSocketReadBufferSize,
		WebSocketWriteBufferSize: DefaultWebSocketWriteBufferSize,
		MaxConnections:           DefaultMaxConnections,
		HeartbeatInterval:        DefaultHeartbeatInterval,
		SessionTimeout:           DefaultSessionTimeout,
	}
}

// Load reads configuration from environment variables and returns a Config instance.
// Missing environment variables will use sensible defaults.
// Returns an error if any required validation fails.
func Load() (*Config, error) {
	config := defaultConfig()

	// Load environment variables with type conversion
	if err := loadEnvInt("PORT", &config.Port); err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	loadEnvString("HOST", &config.Host)

	loadEnvString("CORS_ORIGIN", &config.CORSOrigin)

	loadEnvString("LOG_LEVEL", &config.LogLevel)

	loadEnvString("ENV", &config.Environment)

	if err := loadEnvInt("WS_READ_BUFFER_SIZE", &config.WebSocketReadBufferSize); err != nil {
		return nil, fmt.Errorf("invalid WS_READ_BUFFER_SIZE: %w", err)
	}

	if err := loadEnvInt("WS_WRITE_BUFFER_SIZE", &config.WebSocketWriteBufferSize); err != nil {
		return nil, fmt.Errorf("invalid WS_WRITE_BUFFER_SIZE: %w", err)
	}

	if err := loadEnvInt("MAX_CONNECTIONS", &config.MaxConnections); err != nil {
		return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %w", err)
	}

	if err := loadEnvInt("HEARTBEAT_INTERVAL", &config.HeartbeatInterval); err != nil {
		return nil, fmt.Errorf("invalid HEARTBEAT_INTERVAL: %w", err)
	}

	if err := loadEnvInt("SESSION_TIMEOUT", &config.SessionTimeout); err != nil {
		return nil, fmt.Errorf("invalid SESSION_TIMEOUT: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate checks that the configuration values are valid.
// Returns an error if any configuration value is invalid.
func (c *Config) Validate() error {
	if err := c.validateBasicFields(); err != nil {
		return err
	}

	if err := c.validateLogLevel(); err != nil {
		return err
	}

	if err := c.validateEnvironment(); err != nil {
		return err
	}

	if err := c.validateWebSocketSettings(); err != nil {
		return err
	}

	if err := c.validateConnectionSettings(); err != nil {
		return err
	}

	return nil
}

// validateBasicFields validates basic configuration fields.
func (c *Config) validateBasicFields() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	return nil
}

// validateLogLevel validates the log level configuration.
func (c *Config) validateLogLevel() error {
	if c.LogLevel == "" {
		return fmt.Errorf("log level cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(c.LogLevel)] {
		return fmt.Errorf("invalid log level %q, must be one of: debug, info, warn, error", c.LogLevel)
	}

	return nil
}

// validateEnvironment validates the environment configuration.
func (c *Config) validateEnvironment() error {
	if c.Environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}

	validEnvironments := map[string]bool{
		"development": true,
		"production":  true,
		"test":        true,
	}
	if !validEnvironments[strings.ToLower(c.Environment)] {
		return fmt.Errorf("invalid environment %q, must be one of: development, production, test", c.Environment)
	}

	return nil
}

// validateWebSocketSettings validates WebSocket-related configuration.
func (c *Config) validateWebSocketSettings() error {
	if c.WebSocketReadBufferSize <= 0 {
		return fmt.Errorf("WebSocket read buffer size must be positive, got %d", c.WebSocketReadBufferSize)
	}

	if c.WebSocketWriteBufferSize <= 0 {
		return fmt.Errorf("WebSocket write buffer size must be positive, got %d", c.WebSocketWriteBufferSize)
	}

	return nil
}

// validateConnectionSettings validates connection-related configuration.
func (c *Config) validateConnectionSettings() error {
	if c.MaxConnections <= 0 {
		return fmt.Errorf("max connections must be positive, got %d", c.MaxConnections)
	}

	if c.HeartbeatInterval <= 0 {
		return fmt.Errorf("heartbeat interval must be positive, got %d", c.HeartbeatInterval)
	}

	if c.SessionTimeout <= 0 {
		return fmt.Errorf("session timeout must be positive, got %d", c.SessionTimeout)
	}

	return nil
}

// IsDevelopment returns true if the current environment is development.
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

// IsProduction returns true if the current environment is production.
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

// IsTest returns true if the current environment is test.
func (c *Config) IsTest() bool {
	return strings.ToLower(c.Environment) == "test"
}

// LogLevelSlog returns the slog.Level corresponding to the configured log level.
func (c *Config) LogLevelSlog() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Address returns the complete server address in the format "host:port".
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// loadEnvString loads a string environment variable into the target pointer.
// If the environment variable is not set, the target value remains unchanged.
func loadEnvString(envVar string, target *string) {
	if value := os.Getenv(envVar); value != "" {
		*target = value
	}
}

// loadEnvInt loads an integer environment variable into the target pointer.
// If the environment variable is not set, the target value remains unchanged.
// Returns an error if the environment variable is set but cannot be parsed as an integer.
func loadEnvInt(envVar string, target *int) error {
	value := os.Getenv(envVar)
	if value == "" {
		return nil // Keep default value
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("cannot parse %s as integer: %w", envVar, err)
	}

	*target = parsed
	return nil
}
