package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	log  *slog.Logger
	addr string
}

func New(addr string, log *slog.Logger) *Server {
	return &Server{
		addr: addr,
		log:  log,
	}
}

func (s *Server) Run() error {
	s.log.Info("starting HTTP server", "addr", s.addr)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return http.ListenAndServe(s.addr, mux)
}
