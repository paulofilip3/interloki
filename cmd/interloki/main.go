package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/paulofilip3/interloki/internal/app"
	"github.com/paulofilip3/interloki/internal/config"
	"github.com/paulofilip3/interloki/internal/source"
)

var rootCmd = &cobra.Command{
	Use:   "interloki",
	Short: "A real-time log viewer with a web UI",
	Long:  "interloki pipes log streams through a processing pipeline and serves them to a browser-based viewer via WebSocket.",
}

var stdinCmd = &cobra.Command{
	Use:   "stdin",
	Short: "Read log lines from stdin",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromFlags(cmd)
		src := source.NewStdinSource(os.Stdin)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		application, err := app.New(cfg, src)
		if err != nil {
			return err
		}
		return application.Run(ctx)
	},
}

var followCmd = &cobra.Command{
	Use:   "follow",
	Short: "Follow one or more log files (like tail -f)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromFlags(cmd)

		files, _ := cmd.Flags().GetStringSlice("file")
		cfg.FilePaths = files

		src := source.NewFileSource(files)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		application, err := app.New(cfg, src)
		if err != nil {
			return err
		}
		return application.Run(ctx)
	},
}

var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: "Accept log lines over TCP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromFlags(cmd)

		addr := getStringFlag(cmd, "listen", "INTERLOKI_SOCKET_ADDR", cfg.SocketAddr)
		cfg.SocketAddr = addr

		src := source.NewSocketSource(addr)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		application, err := app.New(cfg, src)
		if err != nil {
			return err
		}
		return application.Run(ctx)
	},
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Generate fake log messages for demonstration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromFlags(cmd)

		rate := getIntFlag(cmd, "rate", "INTERLOKI_DEMO_RATE", cfg.DemoRate)
		cfg.DemoRate = rate

		src := source.NewDemoSource(rate)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		application, err := app.New(cfg, src)
		if err != nil {
			return err
		}
		return application.Run(ctx)
	},
}

func init() {
	// Persistent flags on root command.
	rootCmd.PersistentFlags().Int("port", 8080, "HTTP server port")
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "HTTP server bind address")
	rootCmd.PersistentFlags().Int("max-messages", 10000, "Ring buffer capacity")
	rootCmd.PersistentFlags().Int("bulk-window-ms", 100, "WebSocket flush interval in milliseconds")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose (debug) logging")

	// follow command flags.
	followCmd.Flags().StringSlice("file", nil, "File paths to follow (can be repeated)")
	followCmd.MarkFlagRequired("file")

	// socket command flags.
	socketCmd.Flags().String("listen", ":9999", "TCP address to listen on")

	// demo command flags.
	demoCmd.Flags().Int("rate", 10, "Messages per second")

	rootCmd.AddCommand(stdinCmd, followCmd, socketCmd, demoCmd)
}

// configFromFlags builds a Config from cobra command flags with env var fallback.
func configFromFlags(cmd *cobra.Command) config.Config {
	cfg := config.DefaultConfig()

	cfg.Port = getIntFlag(cmd, "port", "INTERLOKI_PORT", cfg.Port)
	cfg.Host = getStringFlag(cmd, "host", "INTERLOKI_HOST", cfg.Host)
	cfg.MaxMessages = getIntFlag(cmd, "max-messages", "INTERLOKI_MAX_MESSAGES", cfg.MaxMessages)
	cfg.BulkWindowMS = getIntFlag(cmd, "bulk-window-ms", "INTERLOKI_BULK_WINDOW_MS", cfg.BulkWindowMS)
	cfg.Verbose = getBoolFlag(cmd, "verbose", "INTERLOKI_VERBOSE", cfg.Verbose)

	return cfg
}

// getIntFlag returns the flag value if changed, else the env var if set, else the fallback.
func getIntFlag(cmd *cobra.Command, flag, envVar string, fallback int) int {
	if cmd.Flags().Changed(flag) {
		v, _ := cmd.Flags().GetInt(flag)
		return v
	}
	if env := os.Getenv(envVar); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			return v
		}
	}
	return fallback
}

// getStringFlag returns the flag value if changed, else the env var if set, else the fallback.
func getStringFlag(cmd *cobra.Command, flag, envVar string, fallback string) string {
	if cmd.Flags().Changed(flag) {
		v, _ := cmd.Flags().GetString(flag)
		return v
	}
	if env := os.Getenv(envVar); env != "" {
		return env
	}
	return fallback
}

// getBoolFlag returns the flag value if changed, else the env var if set, else the fallback.
func getBoolFlag(cmd *cobra.Command, flag, envVar string, fallback bool) bool {
	if cmd.Flags().Changed(flag) {
		v, _ := cmd.Flags().GetBool(flag)
		return v
	}
	if env := os.Getenv(envVar); env != "" {
		v, err := strconv.ParseBool(env)
		if err == nil {
			return v
		}
	}
	return fallback
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
