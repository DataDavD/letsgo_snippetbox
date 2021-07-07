package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials error is used if a user tries
	// to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail error is used if a user tries to signup
	// with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID       int
	Name     string
	Email    string
	hashedPw []byte
	Created  time.Time
	Active   bool
}
