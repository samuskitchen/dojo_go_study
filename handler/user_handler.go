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

// NewUserHandler ...
func NewUserHandler(db *database.Data) *UserRouter {
	return &UserRouter{
		Repo: userRepo.NewPostgresUserRepo(db),
	}
}

// UserRouter ...
type UserRouter struct {
	Repo repository.UserRepository
}

// GetAllUser response all the users.
func (ur *UserRouter) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := ur.Repo.GetAllUser(ctx)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	middleware.JSON(w, r, http.StatusOK, users)
}

// GetOneHandler response one user by id.
func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	userResult, err := ur.Repo.GetOne(ctx, uint(id))
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	middleware.JSON(w, r, http.StatusOK, userResult)
}

// CreateHandler Create a new user.
func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	if err := user.HashPassword(); err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repo.Create(ctx, &user)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	user.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), user.ID))
	middleware.JSON(w, r, http.StatusCreated, user)

}

// UpdateHandler update a stored user by id.
func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var userUpdate model.User
	err = json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	ctx := r.Context()

	err = ur.Repo.Update(ctx, uint(id), userUpdate)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	middleware.JSON(w, r, http.StatusOK, nil)
}

// DeleteHandler Remove a user by ID.
func (ur *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repo.Delete(ctx, uint(id))
	if err != nil {
		middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	middleware.JSON(w, r, http.StatusNoContent, middleware.Map{})
}

