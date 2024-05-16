package users

import (
	"api/models"
	"api/utils"
	"fmt"
	"net/http"
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LOGIN")
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
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

func (h *Handler) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	body := struct {
		User_id int64 `json:"user_id"`
	}{}
	var err error

	if err = utils.ValidateBody(&body, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	data, err := GetUserPosts(body.User_id, h.db)

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		data,
		"USER ENCONTRADO",
		http.StatusOK,
		w,
	)
}

func (h *Handler) handleGetAllUsersPosts(w http.ResponseWriter, r *http.Request) {
	data, err := GetAllUserPosts(h.db)

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		data,
		"RESULT",
		http.StatusOK,
		w,
	)
}
