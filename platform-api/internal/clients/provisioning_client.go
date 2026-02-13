package clients

import (
	"time"
)

// RetryConfig for gRCP client with resilience
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// ProvisioningClient wraps gRPC client with resilience
