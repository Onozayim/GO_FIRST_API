package posts

import (
	"api/auth"
	"database/sql"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /post/get", h.handeGetPost)

	router.HandleFunc("POST /post/create", auth.CheckAuth(h.handleCreatePost))
}
