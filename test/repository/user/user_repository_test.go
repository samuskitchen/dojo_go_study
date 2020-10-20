package user

import (
	"context"
	"database/sql"
	"dojo_go_study/config/database"
	"dojo_go_study/model"
	"dojo_go_study/repository"
	repo "dojo_go_study/repository/user"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// represent the repository
var (
	dbMockUsers        *sql.DB
	connMockUser       *database.Data
	userRepositoryMock repository.UserRepository
)

// NewMockUser initialize mock connection to database
func NewMockUser() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMockUsers = db
	connMockUser = &database.Data{
		DB: dbMockUsers,
	}

	userRepositoryMock = repo.NewPostgresUserRepo(connMockUser)

	/*userRepositoryMock = &repo.UserRepository{
		Data: &connMockUser,
	}*/

	return mock
}

// Close attaches the provider and close the connection
func CloseMockUser() {
	err := dbMockUsers.Close()
	if err != nil {
		log.Println("Error close database test")
	}
}

// dataUSer is dataDB for test
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
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func TestUserRepository_GetAll(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		mock.ExpectQuery("SELECT 1 FROM user")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := userRepositoryMock.GetAllUser(ctx)
		assert.Error(tt, err)
		assert.Nil(tt, users)
	})

	t.Run("Get All User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		usersData := dataUSer()
		rows := sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "created_at", "updated_at"}).
			AddRow(usersData[0].ID, usersData[0].Name, usersData[0].Surname, usersData[0].Username, usersData[0].Email, usersData[0].CreatedAt, usersData[0].UpdatedAt).
			AddRow(usersData[1].ID, usersData[1].Name, usersData[1].Surname, usersData[1].Username, usersData[1].Email, usersData[1].CreatedAt, usersData[1].UpdatedAt)

		mock.ExpectQuery(selectAllUserTest).WillReturnRows(rows)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := userRepositoryMock.GetAllUser(ctx)
		assert.NotEmpty(tt, users)
		assert.NoError(tt, err)
		assert.Len(tt, users, 2)
	})
}

func TestUserRepository_GetOne(t *testing.T) {

	usersData := dataUSer()
	userTest := model.User{
		ID:        usersData[0].ID,
		Name:      usersData[0].Name,
		Surname:   usersData[0].Surname,
		Username:  usersData[0].Username,
		Email:     usersData[0].Email,
		Password:  usersData[0].Password,
		CreatedAt: usersData[0].CreatedAt,
		UpdatedAt: usersData[0].UpdatedAt,
	}

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "created_at", "updated_at"}).
			AddRow(userTest.ID, userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUserByIdTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetOne(ctx, 1)
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Get User By Id Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "created_at", "updated_at"}).
			AddRow(userTest.ID, userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUserByIdTest).WithArgs(userTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetOne(ctx, 1)
		assert.NoError(tt, err)
		assert.NotNil(tt, userResult)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {

	usersData := dataUSer()
	userTest := model.User{
		ID:        usersData[0].ID,
		Name:      usersData[0].Name,
		Surname:   usersData[0].Surname,
		Username:  usersData[0].Username,
		Email:     usersData[0].Email,
		Password:  usersData[0].Password,
		CreatedAt: usersData[0].CreatedAt,
		UpdatedAt: usersData[0].UpdatedAt,
	}

	t.Run("Error Scan Row", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"idt", "name", "surname", "username", "email", "password", "created_at", "updated_at"}).
			AddRow(userTest.Name, userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUSerByUsernameTest).WithArgs(userTest.Username).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetByUsername(ctx, "daniel.delapava")
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Get User By Username Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password", "created_at", "updated_at"}).
			AddRow(userTest.ID, userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUSerByUsernameTest).WithArgs(userTest.Username).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetByUsername(ctx, "daniel.delapava")
		assert.NoError(tt, err)
		assert.NotNil(tt, userResult)
	})

}

func TestUserRepository_Create(t *testing.T) {

	usersData := dataUSer()
	userTest := &model.User{
		Name:      usersData[0].Name,
		Surname:   usersData[0].Surname,
		Username:  usersData[0].Username,
		Email:     usersData[0].Email,
		Password:  usersData[0].Password,
		CreatedAt: usersData[0].CreatedAt,
		UpdatedAt: usersData[0].UpdatedAt,
	}

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare("insertUserTest")
		prep.ExpectExec().
			WithArgs(usersData[0].Name, usersData[0].Surname, usersData[0].Username, usersData[0].Email, usersData[0].Password, usersData[0].CreatedAt, usersData[0].UpdatedAt).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Create(ctx, userTest)
		assert.Error(tt, err)

	})

	t.Run("Error Scan Row", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(insertUserTest)
		prep.ExpectQuery().
			WithArgs(userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Error"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Create(ctx, userTest)
		assert.Error(tt, err)
	})

	t.Run("Create User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(insertUserTest)
		prep.ExpectQuery().
			WithArgs(userTest.Name, userTest.Surname, userTest.Username, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Create(ctx, userTest)
		assert.NoError(tt, err)
	})
}

func TestUserRepository_Update(t *testing.T) {

	usersData := dataUSer()
	userTest := model.User{
		ID:        usersData[0].ID,
		Name:      usersData[0].Name,
		Surname:   usersData[0].Surname,
		Username:  usersData[0].Username,
		Email:     usersData[0].Email,
		Password:  usersData[0].Password,
		CreatedAt: usersData[0].CreatedAt,
		UpdatedAt: usersData[0].UpdatedAt,
	}

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare("updateUserTest")
		prep.ExpectExec().
			WithArgs(userTest.Name, userTest.Surname, userTest.Username, userTest.UpdatedAt, userTest.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Update(ctx, 1, userTest)
		assert.Error(tt, err)
	})

	t.Run("Error Statement SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(updateUserTest)
		prep.ExpectExec().
			WithArgs(userTest.Name, userTest.Surname, userTest.Username, userTest.UpdatedAt, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Update(ctx, 1, userTest)
		fmt.Println(err)
		assert.Error(tt, err)
	})

	t.Run("Update User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(updateUserTest)
		prep.ExpectExec().
			WithArgs(userTest.Name, userTest.Surname, userTest.Username, userTest.UpdatedAt, userTest.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Update(ctx, uint(1), userTest)
		assert.NoError(tt, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare("deleteUserTest")
		prep.ExpectExec().
			WithArgs(uint(1)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Delete(ctx, 2)
		assert.Error(tt, err)
	})

	t.Run("Error Statement SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(deleteUserTest)
		prep.ExpectExec().
			WithArgs(nil).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Delete(ctx, 1)
		assert.Error(tt, err)
	})

	t.Run("Delete User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(deleteUserTest)
		prep.ExpectExec().
			WithArgs(uint(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.Delete(ctx, 1)
		assert.NoError(tt, err)
	})
}
