package sql

import (
	"database/sql"
	"log"
)

// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed it alongside the connection pool.
type exampleModelWithStmt struct {
	DB         *sql.DB
	insertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared statement
func NewExampleModel(db *sql.DB) (*exampleModelWithStmt, error) {
	// Use the Prepare method to create a new prepared statement for the
	// current connection pool. This returns a sql.Stmt object which represents
	// the prepared statement.
	inStmt, err := db.Prepare("INSERT INTO...")
	if err != nil {
		return nil, err
	}

	// Store it in our exampleModelWithStmt object, alongside the connection pool
	return &exampleModelWithStmt{db, inStmt}, nil
}

// Any methods implemented against the ExampleModel object will have access to
// the prepared statement.
func (m *exampleModelWithStmt) Insert(args ...) error {
	// Notice how we call Exec directly against the prepared statement, rather
	// than against the connection pool? Prepared statements also support the Query
	// and QueryRow methods.
	if _, err := m.insertStmt.Exec(args...); err != nil {
		log.Printf("prepared Insert statement failed: %v", err)
		return err
	}

	return nil
}

// In the web application's main function we will need to initialize a new
// ExampleModel struct using the constructor function.

func examplMain() {
	db, err := sql.Open(...)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new exampleModelWithStmt object, which includes the prepared statement.
	exampleModelWithStmt, err := NewExampleModel(db)
	if err != nil {
		log.Fatal(err)
	}

	// Defer a call to Close on the prepared statment to ensure that it is properly
	// closed before our main function terminates.
	defer func() {
		err := exampleModelWithStmt.insertStmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
