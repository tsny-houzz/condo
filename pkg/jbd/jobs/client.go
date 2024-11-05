package jobs

import (
	"time"
)

// Constants for default values
const (
	DefaultHost    = "tools.stghouzz.com"
	DefaultScheme  = "https"
	DefaultTimeout = 10 * time.Second
)

// JobClient holds configuration for making requests.
type JobClient struct {
	Host    string
	Scheme  string
	Timeout time.Duration
}

// JobClientOption is a type for setting options on JobClient.
type JobClientOption func(*JobClient)

// WithHost sets a custom host for the JobClient.
func WithHost(host string) JobClientOption {
	return func(jc *JobClient) {
		jc.Host = host
	}
}

// WithScheme sets a custom scheme (http or https) for the JobClient.
func WithScheme(scheme string) JobClientOption {
	return func(jc *JobClient) {
		jc.Scheme = scheme
	}
}

// WithTimeout sets a custom timeout for the JobClient.
func WithTimeout(timeout time.Duration) JobClientOption {
	return func(jc *JobClient) {
		jc.Timeout = timeout
	}
}

// NewJobClient creates a new JobClient with the specified options.
func NewJobClient(opts ...JobClientOption) *JobClient {
	// Initialize with default values
	client := &JobClient{
		Host:    DefaultHost,
		Scheme:  DefaultScheme,
		Timeout: DefaultTimeout,
	}

	// Apply each option
	for _, opt := range opts {
		opt(client)
	}

	return client
}
