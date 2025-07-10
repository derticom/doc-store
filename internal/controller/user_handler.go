package controller

import (
	"encoding/json"
	"net/http"

	"github.com/derticom/doc-store/internal/domain/user"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	uc user.UseCase
}

func NewUserHandler(uc user.UseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Login    string `json:"login"`
		Password string `json:"pswd"`
		Token    string `json:"token"` // admin token
	}
	json.NewDecoder(r.Body).Decode(&in)

	err := h.uc.Register(r.Context(), in.Token, in.Login, in.Password)
	if err != nil {
		writeError(w, 403, err)
		return
	}
	writeJSON(w, 200, map[string]any{"response": map[string]string{"login": in.Login}})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Login    string `json:"login"`
		Password string `json:"pswd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, err)
		return
	}

	token, err := h.uc.Authenticate(r.Context(), input.Login, input.Password)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	writeJSON(w, 200, map[string]any{"response": map[string]string{"token": token}})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	h.uc.Logout(r.Context(), token)
	writeJSON(w, 200, map[string]any{"response": map[string]bool{token: true}})
}
