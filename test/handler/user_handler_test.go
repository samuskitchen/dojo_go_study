package handler

import (
	"bytes"
	"context"
	"dojo_go_study/handler"
	"dojo_go_study/model"
	mockLocal "dojo_go_study/repository/mocks"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// dataUSer is data for test
func dataUSer() []model.User {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []model.User{
		{
			ID:        uint(1),
			Name:      "Daniel",
			Surname:   "De La Pava Suarez",
			Username:  "daniel.delapava",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uint(1),
			Name:      "Rebecca",
			Surname:   "Romero",
			Username:  "rebecca.romero",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func dataMockCreate() *model.User {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return &model.User{
		ID:        uint(1),
		Name:      "Daniel",
		Surname:   "De La Pava Suarez",
		Username:  "daniel.delapava",
		Email:     "daniel.delapava@jikkosoft.com",
		Password:  "123456",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestUserRouter_CreateHandler(t *testing.T) {

	t.Run("Error Body Create Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Create Handler", func(tt *testing.T) {
		dataMockCreate().ID = uint(9999999999)

		marshal, err := json.Marshal(dataMockCreate())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Create", mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Create Handler", func(tt *testing.T) {

		var userTest = dataMockCreate()
		userTest.Username = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})

	t.Run("Create Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataMockCreate())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Create", mock.Anything, mock.Anything).Return(nil)

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})

}

func TestUserRouter_GetAllUser(t *testing.T) {

	t.Run("Error Get All User Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("GetAllUser", mock.Anything, mock.Anything).Return(nil, errors.New("error trace test"))

		testUserHandler.GetAllUsersHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})

	t.Run("Get All User Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		response := httptest.NewRecorder()
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("GetAllUser", mock.Anything).Return(dataUSer(), nil)

		testUserHandler.GetAllUsersHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})

}

func TestUserRouter_GetOneHandler(t *testing.T) {

	t.Run("Error Param Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		mockRepository := &mockLocal.UserRepository{}
		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("GetOne", mock.Anything, mock.Anything).Return(model.User{}, errors.New("error sql")).Once()

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("GetOne", mock.Anything, mock.Anything).Return(model.User{}, nil).Once()

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}

func TestUserRouter_UpdateHandler(t *testing.T) {

	t.Run("Error Param Update Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", nil)

		mockRepository := &mockLocal.UserRepository{}
		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Body Update Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(nil))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))

		mockRepository := &mockLocal.UserRepository{}
		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataMockCreate())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Update Handler", func(tt *testing.T) {

		var userTest = dataMockCreate()
		userTest.Username = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataMockCreate())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestUserRouter_DeleteHandler(t *testing.T) {

	t.Run("Error Param Delete Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/api/v1/users/{id}", nil)

		mockRepository := &mockLocal.UserRepository{}
		testUserHandler := &handler.UserRouter{Repo: mockRepository}

		testUserHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Delete Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/api/v1/users/{id}", nil)
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testUserHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Delete Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/api/v1/users/{id}", nil)
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &mockLocal.UserRepository{}

		testUserHandler := &handler.UserRouter{Repo: mockRepository}
		mockRepository.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		testUserHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}
