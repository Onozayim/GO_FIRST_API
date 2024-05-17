package users

import (
	"api/auth"
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
	var errors []string

	if err = utils.ValidateBody(&user, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if !utils.ValidateEmail(user.Email) {
		errors = append(errors, "Email no valido")
	}

	pass_errors := utils.ValidatePassword(user.Password)
	if len(pass_errors) > 0 {
		errors = append(errors, pass_errors...)
	}

	if utils.DbExists(
		h.db,
		`SELECT id FROM users as u WHERE u.email = ?`,
		user.Email,
	) {
		errors = append(errors, "Email ya registrado")
	}

	if len(errors) != 0 {
		utils.ReturnErrorArray(errors, http.StatusBadRequest, w)
		return
	}

	hash, err := auth.HashPassword(user.Password)

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	user.Password = hash

	if err = CreateUser(&user, h.db); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]string{"message": "Usuario Creado!"},
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
		http.StatusOK,
		w,
	)
}
