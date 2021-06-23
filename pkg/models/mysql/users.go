package mysql

import (
	"database/sql"

	"github.com/DataDavD/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert a new user record into the snippetbox.users table
func (u *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies a user exists with the provided email address and password.
// If the user exists the relevant user ID is returned.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get fetches details for a specifc user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
