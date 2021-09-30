package mysql

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our setup and teardown
	// scripts contains multiple SQL statements, we need to use the 'multiStatements=true' parameter
	// in our DSN. This instructs our MySQL database driver to support executing multiple SQL
	// statements in one db.Exec() call.
	dbString := "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true"
	db, err := sql.Open("mysql", dbString)
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Return the connection pool and an anonymous function which reads and executs the teardown
	// script, and closes the connection pool. We can assign this anonymous function and call
	// it later once our test has completed.
	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			err := db.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

}
