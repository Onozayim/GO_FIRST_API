package users

import (
	"api/utils"
	"database/sql"
	"fmt"
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
	router.HandleFunc("POST /user/create", h.handleCreateUser)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LOGIN")
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	var err error

	if err = utils.ValidateBody(&user, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if err = CreateUser(&user, h.db); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]int64{"user_id": user.Id},
		"Usuario Creado",
		http.StatusOK,
		w,
	)
}
