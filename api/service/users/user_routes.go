package users

import (
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
	router.HandleFunc("GET /user/login", h.handleLogin)
	router.HandleFunc("GET /user/get_posts", h.handleGetPosts)
	router.HandleFunc("GET /user/get_all_users_posts", h.handleGetAllUsersPosts)

	router.HandleFunc("POST /user/create", h.handleCreateUser)
}
