package sql

import (
	"database/sql"
	"log"
)

// ExampleModel is an example struct for sql transaction handling
type exampleModel struct {
	DB *sql.DB
}

// ExampleTransaction is an example transaction method on the ExampleModel
func (m *exampleModel) exampleTransaction() error {
	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in-progress database transaction.
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// Call Exec() on the transaction, passing in your statement and any
	// parameters. It's important to notice that tx.Exec() is called on the
	// transaction object just created, NOT the connection pool. Although we're
	// using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
	// exactly the same way.
	_, err = tx.Exec("INSERT INTO...")
	if err != nil {
		// If there is any error, we call the tx.Rollback() method on the
		// transaction. This will abort the transaction and no changes will be
		// made to the database.
		if rb := tx.Rollback(); rb != nil {
			log.Printf("query failed: %v, unable to abort: %v", err, rb)
			return rb
		}
		log.Printf("query failed: %v", err)
		return err
	}

	// If there are no errors, the statements in the transaction can be committed
	// to the database with the tx.Commit() method. It's really important to ALWAYS
	// call either Rollback() or Commit() before your function returns. If you
	// don't the connection will stay open and not be returned to the connection
	// pool. This can lead to hitting your maximum connection limit/running out of
	// resources.
	if err := tx.Commit(); err != nil {
		log.Printf("transaction commit failed: %v", err)
		return err
	}

	return nil
}
