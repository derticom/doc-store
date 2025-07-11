package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/derticom/doc-store/internal/domain/document"
	"github.com/derticom/doc-store/internal/middleware"
	"github.com/derticom/doc-store/internal/usecase/auth"

	"github.com/go-chi/chi/v5"
)

type DocumentHandler struct {
	uc        document.UseCase
	authStore auth.SessionStore
}

func NewDocumentHandler(uc document.UseCase, authStore auth.SessionStore) *DocumentHandler {
	return &DocumentHandler{uc: uc, authStore: authStore}
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	docs, err := h.uc.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"docs": docs,
		},
	})
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	login := r.URL.Query().Get("login")

	doc, data, err := h.uc.Get(r.Context(), id, login)
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}

	if doc.File {
		w.Header().Set("Content-Type", doc.Mime)
		w.WriteHeader(http.StatusOK)
		if r.Method == "GET" {
			w.Write(data)
		}
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": data,
		})
	}
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Чтение meta
	metaPart, _, err := r.FormFile("meta")
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("missing meta"))
		return
	}
	defer metaPart.Close()

	var meta struct {
		Name   string   `json:"name"`
		File   bool     `json:"file"`
		Public bool     `json:"public"`
		Token  string   `json:"token"`
		Mime   string   `json:"mime"`
		Grant  []string `json:"grant"`
	}
	if err := json.NewDecoder(metaPart).Decode(&meta); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// Валидация токена через authStore
	userID, err := h.authStore.GetUserID(r.Context(), meta.Token)
	if err != nil {
		writeError(w, http.StatusForbidden, errors.New("invalid token"))
		return
	}

	// Чтение JSON (опционально)
	var jsonData []byte
	jsonPart, _, _ := r.FormFile("json")
	if jsonPart != nil {
		defer jsonPart.Close()
		jsonData, _ = io.ReadAll(jsonPart)
	}

	// Чтение файла (опционально)
	var file []byte
	if meta.File {
		filePart, _, err := r.FormFile("file")
		if err != nil {
			writeError(w, http.StatusBadRequest, errors.New("missing file"))
			return
		}
		defer filePart.Close()
		file, _ = io.ReadAll(filePart)
	}

	doc := &document.Document{
		Name:     meta.Name,
		File:     meta.File,
		Public:   meta.Public,
		Mime:     meta.Mime,
		Grant:    meta.Grant,
		OwnerID:  userID,
		JSONData: jsonData,
	}

	if err := h.uc.Upload(r.Context(), doc, file); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"file": doc.Name,
			"json": doc.JSONData,
		},
	})
}

func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value(middleware.UserIDKey).(string)

	if err := h.uc.Delete(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"response": map[string]bool{id: true}})
}
