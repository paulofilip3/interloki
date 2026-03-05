package app

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/paulofilip3/interloki/internal/buffer"
	"github.com/paulofilip3/interloki/internal/config"
	"github.com/paulofilip3/interloki/internal/models"
	"github.com/paulofilip3/interloki/internal/processing"
	"github.com/paulofilip3/interloki/internal/server"
	"github.com/paulofilip3/interloki/internal/source"
)

// App wires all interloki components together: source -> pipeline -> buffer -> server.
type App struct {
	cfg config.Config
	src source.Source
}

// New creates a new App with the given configuration and log source.
func New(cfg config.Config, src source.Source) (*App, error) {
	if cfg.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	return &App{cfg: cfg, src: src}, nil
}

// Run starts all components and blocks until ctx is cancelled.
func (a *App) Run(ctx context.Context) error {
	// 1. Create ring buffer.
	ring := buffer.NewRing[models.LogMessage](a.cfg.MaxMessages)

	// 2. Create client manager.
	manager := server.NewClientManager(ring, a.cfg.BulkWindowMS)

	// 3. Create processing pipeline.
	pipe, err := processing.NewPipeline()
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	// 4. Start the source.
	log.WithField("source", a.src.Name()).Info("starting source")
	sourceCh, err := a.src.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start source %q: %w", a.src.Name(), err)
	}

	// 5. Run the pipeline.
	out, errs := pipe.Run(ctx, sourceCh)

	// 6. Start ConsumeLoop in a goroutine.
	go manager.ConsumeLoop(ctx, out)

	// 7. Start error drain goroutine.
	go func() {
		for err := range errs {
			log.WithError(err).Warn("pipeline error")
		}
	}()

	// 8. Create and start HTTP server.
	srv := server.NewServer(a.cfg.Host, a.cfg.Port, manager)

	// 9. Print startup message.
	addr := fmt.Sprintf("http://%s:%d", a.cfg.Host, a.cfg.Port)
	fmt.Printf("interloki listening on %s\n", addr)
	log.WithField("addr", addr).Info("server starting")

	// 10. Block until server returns.
	return srv.Start(ctx)
}
