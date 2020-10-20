package integration

import (
	"bytes"
	testDB "dojo_go_study/config/database/test"
	"dojo_go_study/config/middleware"
	"dojo_go_study/model"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func TestIntegration_GetAllUser(t *testing.T) {

	t.Run("No Content (no seed data)", func(tt *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		result := middleware.Response{}
		if err := json.Unmarshal([]byte(w.Body.String()), &result); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if (middleware.Response{}) == result {
			tt.Errorf("expected no result to be returned, got %v result", result)
		}
	})

	t.Run("Ok (database has been seeded)", func(tt *testing.T) {
		defer func() {
			if err := testDB.Truncate(database.DB); err != nil {
				tt.Errorf("error truncating test database tables: %v", err)
			}
		}()

		expectedUsers, err := testDB.SeedUsers(database.DB)
		if err != nil {
			tt.Fatalf("error seeding users: %v", err)
		}

		resultSeed := middleware.Response{
			Status:  true,
			Data:    middleware.Map{"users": expectedUsers},
			Message: "Ok",
		}

		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		result := middleware.Response{}
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(resultSeed.Status, result.Status); d != "" {
			tt.Errorf("unexpected difference in response body:\n%v", d)
		}
	})
}

func TestIntegration_GetOneHandler(t *testing.T) {
	defer func() {
		if err := testDB.Truncate(database.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedLists, err := testDB.SeedUsers(database.DB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		ExpectedBody middleware.Response
		ExpectedCode int
	}{
		{
			Name:         "Get One User Successful",
			UserID:       expectedLists[0].ID,
			ExpectedBody: middleware.Response{},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "User Not Found",
			UserID:       0,
			ExpectedBody: middleware.Response{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", test.UserID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var response middleware.Response

				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := true, response.Status; e != a {
					t.Errorf("expected status: %v, got response status ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_CreateHandler(t *testing.T) {

	defer func() {
		if err := testDB.Truncate(database.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	tests := []struct {
		Name         string
		RequestBody  model.User
		ExpectedCode int
	}{
		{
			Name:         "Create User Successful",
			RequestBody:  dataUSer()[0],
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "Break Unique UserName Constraint",
			RequestBody:  dataUSer()[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			RequestBody:  model.User{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict && test.ExpectedCode != http.StatusUnprocessableEntity {
				var response middleware.Response

				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := true, response.Status; e != a {
					t.Errorf("expected status: %v, got response status: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_UpdateHandler(t *testing.T) {

	defer func() {
		if err := testDB.Truncate(database.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := testDB.SeedUsers(database.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		RequestBody  model.User
		ExpectedCode int
	}{
		{
			Name:         "Update User Successful",
			UserID:       expectedUsers[0].ID,
			RequestBody:  expectedUsers[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "Break Unique UserName Constraint",
			UserID:       expectedUsers[1].ID,
			RequestBody:  expectedUsers[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			UserID:       expectedUsers[0].ID,
			RequestBody:  model.User{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d", test.UserID), &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_DeleteHandler(t *testing.T) {
	defer func() {
		if err := testDB.Truncate(database.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedLists, err := testDB.SeedUsers(database.DB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		ExpectedCode int
	}{
		{
			Name:         "Delete User Successful",
			UserID:       expectedLists[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", test.UserID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}
