// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	Server      struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	JWT struct {
		Secret string        `yaml:"secret"`
		Expiry time.Duration `yaml:"expiry"`
	} `yaml:"jwt"`
	Redis struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	Plaid struct {
		ClientID string `yaml:"client_id"`
		Secret   string `yaml:"secret"`
	} `yaml:"plaid"`
	OpenAI struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"openai"`
	Google struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	} `yaml:"google"`
}

// Load loads config from YAML file with environment variable overrides
func Load() (*Config, error) {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment")
	}

	// Default config
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Determine config file path
	configFile := "configs/config.yaml"
	if cfg.Environment == "production" {
		configFile = "configs/config.prod.yaml"
	} else if cfg.Environment == "test" {
		configFile = "configs/config.test.yaml"
	}

	// Read and parse YAML
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Warning: cannot read config file %s: %v\n", configFile, err)
	} else {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Apply environment variable overrides
	cfg.applyEnvOverrides()

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	fmt.Printf("Config loaded from %s (env: %s)\n", configFile, cfg.Environment)
	return cfg, nil
}

func (c *Config) applyEnvOverrides() {
	// Server
	if addr := getEnv("SERVER_ADDRESS", ""); addr != "" {
		c.Server.Address = addr
	}

	// Database
	if url := getEnv("DATABASE_URL", ""); url != "" {
		c.Database.URL = url
	}

	// JWT
	if secret := getEnv("JWT_SECRET", ""); secret != "" {
		c.JWT.Secret = secret
	}
	if expiry := getEnv("JWT_EXPIRY", ""); expiry != "" {
		if d, err := time.ParseDuration(expiry); err == nil {
			c.JWT.Expiry = d
		}
	}

	// Redis
	if url := getEnv("REDIS_URL", ""); url != "" {
		c.Redis.URL = url
	}

	// Plaid
	if id := getEnv("PLAID_CLIENT_ID", ""); id != "" {
		c.Plaid.ClientID = id
	}
	if secret := getEnv("PLAID_SECRET", ""); secret != "" {
		c.Plaid.Secret = secret
	}

	// OpenAI
	if key := getEnv("OPENAI_API_KEY", ""); key != "" {
		c.OpenAI.APIKey = key
	}

	// Google OAuth
	if id := getEnv("GOOGLE_OAUTH_CLIENT_ID", ""); id != "" {
		c.Google.ClientID = id
	}
	if secret := getEnv("GOOGLE_OAUTH_CLIENT_SECRET", ""); secret != "" {
		c.Google.ClientSecret = secret
	}
	if redirect := getEnv("GOOGLE_OAUTH_REDIRECT_URL", ""); redirect != "" {
		c.Google.RedirectURL = redirect
	}
}

func (c *Config) validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Optional but warn if missing
	if c.Plaid.ClientID != "" && c.Plaid.Secret == "" {
		return fmt.Errorf("PLAID_SECRET required if PLAID_CLIENT_ID is set")
	}
	if c.OpenAI.APIKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY not set - AI features will be disabled")
	}
	if c.Google.ClientID == "" || c.Google.ClientSecret == "" {
		fmt.Println("Warning: Google OAuth credentials missing - Google login disabled")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
