package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/DataDavD/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Insert a new user record into the snippetbox.users table
func (u *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec(0 method to insert the user details and hashed password
	// into the users table.
	_, err = u.DB.Exec(stmt, name, email, string(hashedPw))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by checking
		// the contents of the message string. If it does, we return an
		// ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message,
				"users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil

}

// Authenticate verifies a user exists with the provided email address and password.
// If the user exists the relevant user ID is returned.
func (u *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password assocaited with the given email.
	// If no matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPw []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`
	row := u.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.

	err = bcrypt.CompareHashAndPassword(hashedPw, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct, so return the userID.
	return id, nil
}

// Get fetches details for a specific user based on their user ID.
func (u *UserModel) Get(id int) (*models.User, error) {
	usr := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	err := u.DB.QueryRow(stmt, id).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Created, &usr.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return usr, nil
}
