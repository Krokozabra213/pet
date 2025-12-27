package apphttp

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
)

const timeout = 10 * time.Second

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *platformconfig.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) MustRun() {
	if err := s.RunHTTP(); err != nil {
		panic(err)
	}
}

func (s *Server) RunHTTP() error {
	slog.Info("starting application")
	return s.httpServer.ListenAndServe()
}

func (s *Server) StopHTTP() error {
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	return s.httpServer.Shutdown(ctx)
}
