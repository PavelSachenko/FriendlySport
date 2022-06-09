package server

import (
	"fmt"
	"log"
	"net/http"
	"user_service/config"
)

type Server struct {
	httpServer *http.Server
}

func InitServer(cfg *config.Config, handler http.Handler) *Server {
	log.Printf("Initial http server by port: %s", cfg.Server.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}
