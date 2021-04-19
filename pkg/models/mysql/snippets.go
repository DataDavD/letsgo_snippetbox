package mysql

import (
	"database/sql"
	"errors"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.02 (Project Setup and Enabling
	// Modules) so that the import statement looks like this:
	// "{your-module-path}/pkg/models".

	"github.com/DataDavD/snippetbox/pkg/models"
)

// SnippetModel is type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the database. It returns the ID inserted
// and error. If there is no error then Insert returns ID and nil. If there is an error,
// it returns 0 and error.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Write the SQL statement we want to execute. It's split over two lines which
	// why its surrounded with backquotes instead of normal double quotes.
	stmt := `INSERT INTO snippetbox.snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the embedded database connection pool to execute the statement.
	// The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use LastInsertID() method on the result object to get the ID of our newly
	// inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to the int type before returning
	return int(id), nil
}

// Get returns a specific snippet based on the id. It returns ID and error.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippetbox.snippets
	WHERE expires > UTC_TIMESTAMP() and id = ?`

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number  of
	// columns returned by your statement
	if err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own models.ErrNoRecord error instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK Then return the Snippet object.
	return s, nil

}

// Latest returns the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
