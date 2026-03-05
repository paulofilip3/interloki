package config

// Config holds the runtime configuration for interloki.
type Config struct {
	// Server settings
	Host         string
	Port         int
	MaxMessages  int  // ring buffer capacity
	BulkWindowMS int  // WebSocket flush interval in milliseconds
	Verbose      bool

	// Source-specific settings
	FilePaths  []string // for follow command
	SocketAddr  string // for socket command (e.g., ":9999")
	ForwardAddr string // for forward command (e.g., ":24224")
	DemoRate    int    // messages per second for demo command
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Host:         "0.0.0.0",
		Port:         8080,
		MaxMessages:  10000,
		BulkWindowMS: 100,
		Verbose:      false,

		FilePaths:   nil,
		SocketAddr:  ":9999",
		ForwardAddr: ":24224",
		DemoRate:    10,
	}
}
