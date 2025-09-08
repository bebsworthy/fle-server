package config_test

import (
	"os"
	"testing"

	"github.com/fle/server/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	// Clear environment to get defaults
	os.Clearenv()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected default port to be 8080, got %d", cfg.Port)
	}

	if cfg.Host != "0.0.0.0" {
		t.Errorf("Expected default host to be '0.0.0.0', got %s", cfg.Host)
	}

	const expectedEnv = "development"
	if cfg.Environment != expectedEnv {
		t.Errorf("Expected default environment to be '%s', got %s", expectedEnv, cfg.Environment)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got %s", cfg.LogLevel)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear environment to test defaults
	os.Clearenv()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error loading defaults, got: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected default port to be 8080, got %d", cfg.Port)
	}

	const expectedEnv = "development"
	if cfg.Environment != expectedEnv {
		t.Errorf("Expected default environment to be '%s', got %s", expectedEnv, cfg.Environment)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	if err := os.Setenv("PORT", "9000"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("HOST", "localhost"); err != nil {
		t.Fatalf("Failed to set HOST: %v", err)
	}
	if err := os.Setenv("ENV", "production"); err != nil {
		t.Fatalf("Failed to set ENV: %v", err)
	}
	if err := os.Setenv("LOG_LEVEL", "error"); err != nil {
		t.Fatalf("Failed to set LOG_LEVEL: %v", err)
	}

	defer func() {
		_ = os.Unsetenv("PORT")      // Errors are ignored in cleanup
		_ = os.Unsetenv("HOST")      // Errors are ignored in cleanup
		_ = os.Unsetenv("ENV")       // Errors are ignored in cleanup
		_ = os.Unsetenv("LOG_LEVEL") // Errors are ignored in cleanup
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error loading from env, got: %v", err)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected port to be 9000, got %d", cfg.Port)
	}

	if cfg.Host != "localhost" {
		t.Errorf("Expected host to be 'localhost', got %s", cfg.Host)
	}

	const expectedEnv = "production"
	if cfg.Environment != expectedEnv {
		t.Errorf("Expected environment to be '%s', got %s", expectedEnv, cfg.Environment)
	}

	if cfg.LogLevel != "error" {
		t.Errorf("Expected log level to be 'error', got %s", cfg.LogLevel)
	}
}

func TestValidation(t *testing.T) {
	// Clear environment to get defaults
	os.Clearenv()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test valid config
	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config to pass validation, got: %v", err)
	}

	// Test invalid port
	cfg.Port = 0
	if err := cfg.Validate(); err == nil {
		t.Error("Expected invalid port to fail validation")
	}

	// Reset and test invalid log level
	cfg, _ = config.Load()
	cfg.LogLevel = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("Expected invalid log level to fail validation")
	}

	// Reset and test invalid environment
	cfg, _ = config.Load()
	cfg.Environment = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("Expected invalid environment to fail validation")
	}
}

func TestHelperMethods(t *testing.T) {
	// Clear environment to get defaults
	os.Clearenv()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test environment checks
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for default config")
	}

	cfg.Environment = "production"
	if !cfg.IsProduction() {
		t.Error("Expected IsProduction() to return true for production config")
	}

	cfg.Environment = "test"
	if !cfg.IsTest() {
		t.Error("Expected IsTest() to return true for test config")
	}

	// Test address
	cfg, _ = config.Load()
	expected := "0.0.0.0:8080"
	if cfg.Address() != expected {
		t.Errorf("Expected address to be %s, got %s", expected, cfg.Address())
	}
}
