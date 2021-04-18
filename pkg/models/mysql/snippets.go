package mysql

import (
	"database/sql"
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

// Insert inserts a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get returns a specific snippet based on the id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest returns the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
