package http

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewServer(handler http.Handler, config Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":" + config.Port,
			Handler:      handler,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
