package server

import (
	"log/slog"
	"net/http"

	"github.com/derticom/doc-store/internal/controller"
	"github.com/derticom/doc-store/internal/domain/document"
	"github.com/derticom/doc-store/internal/domain/user"
	mddlwr "github.com/derticom/doc-store/internal/middleware"
	"github.com/derticom/doc-store/internal/usecase/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	log       *slog.Logger
	addr      string
	documents document.UseCase
	users     user.UseCase
	authStore auth.SessionStore
}

func New(
	addr string,
	log *slog.Logger,
	documents document.UseCase,
	users user.UseCase,
	authStore auth.SessionStore,
) *Server {
	return &Server{
		addr:      addr,
		log:       log,
		documents: documents,
		users:     users,
		authStore: authStore,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userHandler := controller.NewUserHandler(s.users)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/auth", userHandler.Login)
		r.Delete("/auth/{token}", userHandler.Logout)
	})

	docHandler := controller.NewDocumentHandler(s.documents, s.authStore)

	r.Route("/api/docs", func(r chi.Router) {
		r.Post("/", docHandler.Upload)

		r.Group(func(r chi.Router) {
			r.Use(mddlwr.AuthMiddleware(s.authStore))

			r.Get("/", docHandler.List)
			r.Head("/", docHandler.List)

			r.Get("/{id}", docHandler.Get)
			r.Head("/{id}", docHandler.Get)

			r.Delete("/{id}", docHandler.Delete)
		})
	})

	s.log.Info("Starting server", "addr", s.addr)
	return http.ListenAndServe(s.addr, r)
}
