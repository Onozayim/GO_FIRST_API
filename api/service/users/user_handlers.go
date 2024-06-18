package users

import (
	"api/auth"
	"api/models"
	"api/utils"
	"fmt"
	"net/http"
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	user2 := &models.User{}
	var err error

	if err = utils.ValidateBody(&user, w, r); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if !utils.ValidateEmail(user.Email) {
		utils.ReturnErrorStatus(fmt.Errorf("email no valido"), http.StatusBadRequest, w)
		return
	}

	if user2, err = GetUserByEmail(h.db, user.Email); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if !auth.CheckPassword(user.Password, user2.Password) {
		utils.ReturnErrorStatus(fmt.Errorf("las credenciales no coinciden"), http.StatusBadRequest, w)
		return
	}

	token, err := auth.CreateToken(*user2)

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]string{"message": "Usuario Logeado!", "token": token},
		http.StatusOK,
		w,
	)
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	var err error
	var errors []string

	if err = utils.ValidateBody(user, w, r); err != nil {
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

	hash, err := auth.HashPassword(user.Password)

	if err != nil {
		errors = append(errors, err.Error())
	}

	user.Password = hash

	if err = CreateUser(&user, h.db); err != nil {
		errors = append(errors, err.Error())
	}

	token, err := auth.CreateToken(user)

	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) != 0 {
		utils.ReturnErrorArray(errors, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]string{"message": "Usuario Creado!", "token": token},
		http.StatusOK,
		w,
	)
}

func (h *Handler) handleGetPostsWithJWT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	user := models.User{}

	if err = auth.GetUserNameFromContext(ctx, &user); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	data, err := GetUserPosts(user.Id, h.db)

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

func (h *Handler) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	user := models.User{}

	if err = auth.GetUserNameFromContext(ctx, &user); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if err = utils.SendEmailHTML(
		"Prueba",
		[]string{user.Email},
		"./templates/mails/test.html",
		struct{ Name string }{Name: user.Username},
	); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]string{"status": "done"},
		http.StatusOK,
		w,
	)
}

func (h *Handler) handlePostProfilePicture(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	user := models.User{}

	if err = auth.GetUserNameFromContext(ctx, &user); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	file, handler, err := r.FormFile("picture")

	if err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	if err = utils.SaveFile(file, handler, "", true); err != nil {
		utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
		return
	}

	utils.ReturnOkStatus(
		map[string]string{"status": "done"},
		http.StatusOK,
		w,
	)
}
