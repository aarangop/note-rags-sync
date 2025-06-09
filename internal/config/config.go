package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application config
	Version   string
	VaultPath string

	// AWS config
	S3Bucket  string
	AWSRegion string

	// Optional: Other settings
	LogLevel string
	LogFile  string
	HTTPPort int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Version:   getEnvWithDefault("APP_VERSION", "dev"),
		VaultPath: getEnvWithDefault("VAULT_PATH", ""),
		S3Bucket:  getEnvWithDefault("S3_BUCKET", ""),
		AWSRegion: getEnvWithDefault("AWS_REGION", "us-east-1"),
		LogLevel:  getEnvWithDefault("LOG_LEVEL", "info"),
		LogFile:   getEnvWithDefault("LOG_FILE", "logs/obsidian-sync.log"),
	}

	if portStr := os.Getenv("HTTP_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_PORT: %v", err)
		}
		cfg.HTTPPort = port
	} else {
		cfg.HTTPPort = 8080
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failde: %v", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	// Check required fields
	if c.VaultPath == "" {
		return fmt.Errorf("VAULT_PATH is required")
	}

	if _, err := os.Stat(c.VaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault path does not exist: %s", c.VaultPath)
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// String returns a string representation (useful for logging)
// This implements the Stringer interface we discussed earlier
func (c *Config) String() string {
	return fmt.Sprintf("Config{Version: %s, VaultPath: %s, S3Bucket: %s, AWSRegion: %s, LogLevel: %s}",
		c.Version, c.VaultPath, c.S3Bucket, c.AWSRegion, c.LogLevel)
}

// SetupLogging initializes the logger with configuration from this Config
func (c *Config) SetupLogging() {
	// Import the logger here to avoid import cycle
	// This will be called by main, not by config itself
}
