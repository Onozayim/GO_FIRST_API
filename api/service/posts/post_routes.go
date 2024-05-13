package posts

import (
	"api/utils"
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
	router.HandleFunc("POST /post/create", h.handleCreatePost)
	router.HandleFunc("GET /post/get", h.GetPosts)
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	var err error

	if err = utils.ValidateBody(&post, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if err = CreatePost(&post, h.db); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]int64{"post_id": post.Id},
		"Post creado",
		http.StatusOK,
		w,
	)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Post_id int64 `json:"post_id"`
	}{}
	var err error

	if err = utils.ValidateBody(&body, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	post, err := GetPost(body.Post_id, h.db)

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(post, "POST ENCONTRADO", http.StatusOK, w)
}
