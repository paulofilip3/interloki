package config

// Config holds the runtime configuration for interloki.
type Config struct {
	Host         string
	Port         int
	MaxMessages  int  // ring buffer capacity
	BulkWindowMS int  // WebSocket flush interval in milliseconds
	Verbose      bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Host:         "0.0.0.0",
		Port:         8080,
		MaxMessages:  10000,
		BulkWindowMS: 100,
		Verbose:      false,
	}
}
