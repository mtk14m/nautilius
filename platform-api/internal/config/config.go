package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
// Supports env var overrides

type Config struct {
	// Server
	Server ServerConfig

	// gRPC clients
	ProvisioningService ProvisioningServiceConfig

	// Observability
	Observability ObservabilityConfig

	// Auth
	Auth AuthConfig
}

type ServerConfig struct {
	Port           int
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type ProvisioningServiceConfig struct {
	// gRPC endpoint (e.g., "provisioning-service:50051")
	Host string
	Port int

	// Connection settings
	MaxConnections    int
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration

	// Retry policy
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

type ObservabilityConfig struct {
	LogLevel  string // debug, info, warn, error
	LogFormat string // json, text

	// Tracing (OpenTelemetry)
	TracingEnabled     bool
	TracingServiceName string
	TracingEndpoint    string

	// Metrics
	MetricsEnabled bool
	MetricsPort    int
}

type AuthConfig struct {
	// OIDC provider (e.g., Keycloak)
	Enabled      bool
	IssuerURL    string
	ClientID     string
	ClientSecret string

	// JWT validation
	JWTAudience   string
	TokenCacheTTL time.Duration
}

// Load configuration from env vars with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           getEnvInt("PORT", 8000),
			Addr:           getEnvString("ADDR", "0.0.0.0"),
			ReadTimeout:    getEnvDuration("READ_TIMEOUT", 15*time.Second),
			WriteTimeout:   getEnvDuration("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:    getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
			MaxHeaderBytes: 1 << 20, // 1MB
		},
		ProvisioningService: ProvisioningServiceConfig{
			Host:              getEnvString("PROVISIONING_HOST", "localhost"),
			Port:              getEnvInt("PROVISIONING_PORT", 50051),
			MaxConnections:    getEnvInt("PROVISIONING_MAX_CONN", 100),
			ConnectionTimeout: getEnvDuration("PROVISIONING_CONN_TIMEOUT", 5*time.Second),
			RequestTimeout:    getEnvDuration("PROVISIONING_REQUEST_TIMEOUT", 30*time.Second),
			MaxRetries:        getEnvInt("PROVISIONING_MAX_RETRIES", 3),
			InitialBackoff:    getEnvDuration("PROVISIONING_INITIAL_BACKOFF", 100*time.Millisecond),
			MaxBackoff:        getEnvDuration("PROVISIONING_MAX_BACKOFF", 10*time.Second),
		},
		Observability: ObservabilityConfig{
			LogLevel:           getEnvString("LOG_LEVEL", "info"),
			LogFormat:          getEnvString("LOG_FORMAT", "json"),
			TracingEnabled:     getEnvBool("TRACING_ENABLED", false),
			TracingServiceName: getEnvString("TRACING_SERVICE_NAME", "platform-api"),
			TracingEndpoint:    getEnvString("TRACING_ENDPOINT", "http://localhost:4317"),
			MetricsEnabled:     getEnvBool("METRICS_ENABLED", true),
			MetricsPort:        getEnvInt("METRICS_PORT", 9090),
		},
		Auth: AuthConfig{
			Enabled:       getEnvBool("AUTH_ENABLED", false),
			IssuerURL:     getEnvString("OIDC_ISSUER_URL", ""),
			ClientID:      getEnvString("OIDC_CLIENT_ID", ""),
			ClientSecret:  getEnvString("OIDC_CLIENT_SECRET", ""),
			JWTAudience:   getEnvString("JWT_AUDIENCE", "platform-api"),
			TokenCacheTTL: getEnvDuration("TOKEN_CACHE_TTL", 5*time.Minute),
		},
	}
}

// Helper functions to read env vars with type conversion
func getEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1" || val == "yes"
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Validate checks if config is valid
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.ProvisioningService.Host == "" {
		return fmt.Errorf("provisioning service host is required")
	}

	return nil
}
