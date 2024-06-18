package users

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
	router.HandleFunc("POST /user/login", h.handleLogin)
	router.HandleFunc("GET /user/get_posts", h.handleGetPosts)
	router.HandleFunc("GET /user/get_posts_with_jwt", auth.CheckAuth(h.handleGetPostsWithJWT))
	router.HandleFunc("GET /user/get_all_users_posts", h.handleGetAllUsersPosts)

	router.HandleFunc("POST /user/create", h.handleCreateUser)
	router.HandleFunc("POST /user/email", auth.CheckAuth(h.handleSendEmail))
	router.HandleFunc("POST /user/upload_image", auth.CheckAuth(h.handlePostProfilePicture))
}
