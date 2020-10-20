package handler

import (
	"dojo_go_study/config/database"
	"dojo_go_study/config/middleware"
	"dojo_go_study/model"
	"dojo_go_study/repository"
	userRepo "dojo_go_study/repository/user"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

// UserRouter ...
type UserRouter struct {
	Repo repository.UserRepository
}

// NewUserHandler ...
func NewUserHandler(db *database.Data) *UserRouter {
	return &UserRouter{
		Repo: userRepo.NewPostgresUserRepo(db),
	}
}

// GetAllUser response all the users.
func (ur *UserRouter) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := ur.Repo.GetAllUser(ctx)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound,"0004", err.Error())
		return
	}

	result := middleware.Response{
		Status:  true,
		Data:    middleware.Map{"users": users} ,
		Message: "Ok",
	}

	middleware.JSON(w, r, http.StatusOK, result)
}

// GetOneHandler response one user by id.
func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	ctx := r.Context()
	userResult, err := ur.Repo.GetOne(ctx, uint(id))
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound, "0004", err.Error())
		return
	}

	result := middleware.Response{
		Status:  true,
		Data:    middleware.Map{"user": userResult} ,
		Message: "Ok",
	}

	middleware.JSON(w, r, http.StatusOK, result)
}

// CreateHandler Create a new user.
func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	defer r.Body.Close()
	if err := user.HashPassword(); err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	err = user.Validate("")
	if err != nil {
		middleware.HTTPError(w, r, http.StatusUnprocessableEntity, "0004", err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repo.Create(ctx, &user)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusConflict, "0004", err.Error())
		return
	}

	user.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), user.ID))

	result := middleware.Response{
		Status:  true,
		Data:    middleware.Map{"user": user} ,
		Message: "Ok",
	}

	middleware.JSON(w, r, http.StatusCreated, result)

}

// UpdateHandler update a stored user by id.
func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	var userUpdate model.User
	err = json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	err = userUpdate.Validate("")
	if err != nil {
		middleware.HTTPError(w, r, http.StatusUnprocessableEntity, "0004", err.Error())
		return
	}

	defer r.Body.Close()
	err = ur.Repo.Update(ctx, uint(id), userUpdate)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusConflict, "0004", err.Error())
		return
	}

	result := middleware.Response{
		Status:  true,
		Data:    middleware.Map{"user": userUpdate} ,
		Message: "Ok",
	}

	middleware.JSON(w, r, http.StatusOK, result)
}

// DeleteHandler Remove a user by ID.
func (ur *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, "0004", err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repo.Delete(ctx, uint(id))
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound, "0004", err.Error())
		return
	}

	result := middleware.Response{
		Status:  true,
		Data:    middleware.Map{"id": id} ,
		Message: "Ok",
	}

	middleware.JSON(w, r, http.StatusNoContent, result)
}
