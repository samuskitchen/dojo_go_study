package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uint      `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Surname      string    `json:"surname,omitempty"`
	Username     string    `json:"username,omitempty"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	CreatedAtInt  uint64    `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"-"`
	UpdatedAtInt uint64    `json:"updated_at,omitempty"`
}

// HashPassword generates a hash of the password and places the result in PasswordHash.
func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}

// PasswordMatch compares HashPassword with the password and returns true if they match.
func (u *User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}
