package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPing tests ping handler for the correct response status code, 200 and
// the correct response body, "OK".
func TestPing(t *testing.T) {
	t.Parallel()

	// Create a new instance of our application struct. For now, this just contains
	// a couple of mock loggers (which discard anything written to them).
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	// We then use the httptest.NewTLSServer() function to create a new test server,
	// passing in the value returned by our app.routes() method as the handler for the
	// server. This starts up an HTTPS server which listens on a randomly-chosen port of
	// your local machine for the duration of the test. Notice that we defer a call to
	// ts.Close() to shutdown the server then the test finishes.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained in the
	// ts.URL field. We can use this along with the ts.Client.Get() method to make
	// a GET /healthcheck request against the test server. This returns a http.Response
	// struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Fatal()
	}

	// We can then check the value of the response status code and the body using the
	// same code as before
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// And we can check that the response body written by the ping handler equals "OK".
	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
