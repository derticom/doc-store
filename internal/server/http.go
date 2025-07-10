package server

import (
	"log/slog"
	"net/http"

	"github.com/derticom/doc-store/internal/controller"
	"github.com/derticom/doc-store/internal/domain/document"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	log       *slog.Logger
	addr      string
	documents document.UseCase
}

func New(addr string, log *slog.Logger, documents document.UseCase) *Server {
	return &Server{
		addr:      addr,
		log:       log,
		documents: documents,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Ping endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	docHandler := controller.NewDocumentHandler(s.documents)

	r.Route("/api/docs", func(r chi.Router) {
		r.Get("/", docHandler.List)
		r.Head("/", docHandler.List)
		r.Get("/{id}", docHandler.Get)
		r.Head("/{id}", docHandler.Get)
	})

	s.log.Info("Starting server", "addr", s.addr)
	return http.ListenAndServe(s.addr, r)
}
