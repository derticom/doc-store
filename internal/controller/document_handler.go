package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/derticom/doc-store/internal/domain/document"
)

type DocumentHandler struct {
	uc document.UseCase
}

func NewDocumentHandler(uc document.UseCase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromToken(r) // пока захардкожен

	docs, err := h.uc.List(r.Context(), login)
	if err != nil {
		writeError(w, 500, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"docs": docs,
		},
	})
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/docs/")
	login := getLoginFromToken(r)

	doc, data, err := h.uc.Get(context.Background(), id, login)
	if err != nil {
		writeError(w, 403, err)
		return
	}

	if doc.File {
		w.Header().Set("Content-Type", doc.Mime)
		w.WriteHeader(200)
		if r.Method == "GET" {
			w.Write(data)
		}
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": data,
		})
	}
}
