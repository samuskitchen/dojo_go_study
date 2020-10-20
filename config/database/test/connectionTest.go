package testdb

import (
	"github.com/pkg/errors"
	"time"

	"database/sql"
	db "dojo_go_study/config/database"
	"dojo_go_study/model"
)

// Open returns a new database connection for the test database.
func Open() *db.Data {
	return db.NewTest()
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sql.DB) error {
	stmt := "TRUNCATE TABLE users RESTART IDENTITY CASCADE;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

// SeedUsers handles seeding the user table in the database for integration tests.
func SeedUsers(dbc *sql.DB) ([]model.User, error) {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	users := []model.User{
		{
			Name:      "Daniel",
			Surname:   "De La Pava Suarez",
			Username:  "daniel.delapava",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Rebecca",
			Surname:   "Romero",
			Username:  "rebecca.romero",
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for i := range users {
		query := `INSERT INTO users (
				name, 
				surname, 
				username, 
				email, 
				password, 
				created_at, 
				updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare user insertion")
		}

		row := stmt.QueryRow(&users[i].Name, &users[i].Surname, &users[i].Username, &users[i].Email, &users[i].Password, &users[i].CreatedAt, &users[i].UpdatedAt)

		if err = row.Scan(&users[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture user id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return users, nil
}
