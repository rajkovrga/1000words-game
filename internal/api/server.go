package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"1000words-game/internal/config"
)

type Server struct {
	httpServer *http.Server
	address    string
}

func NewServer(cfg config.Config, handler http.Handler) *Server {
	address := fmt.Sprintf(":%s", cfg.APIPort)

	httpServer := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.APIReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.APIWriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.APIIdleTimeoutSeconds) * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		address:    address,
	}
}

func (s *Server) Start() error {
	fmt.Printf("API server started on http://localhost%s\n", s.address)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
