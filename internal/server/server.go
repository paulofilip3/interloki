package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server is the HTTP/WebSocket server for interloki.
type Server struct {
	httpServer *http.Server
	manager    *ClientManager
}

// NewServer creates a new Server bound to the given host and port.
func NewServer(host string, port int, manager *ClientManager) *Server {
	mux := http.NewServeMux()

	mux.Handle("GET /", FrontendHandler())

	mux.HandleFunc("GET /ws", manager.HandleWS)
	mux.HandleFunc("GET /api/status", manager.HandleStatus)
	mux.HandleFunc("GET /api/client/load", manager.HandleLoadRange)

	return &Server{
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(host, fmt.Sprintf("%d", port)),
			Handler: mux,
		},
		manager: manager,
	}
}

// Handler returns the HTTP handler (mux) for use in tests.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

// Start runs the HTTP server and blocks until ctx is cancelled, at which
// point it initiates a graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
