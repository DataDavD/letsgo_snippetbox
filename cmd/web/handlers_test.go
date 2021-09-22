package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

// TestPing tests ping handler for the correct response status code, 200 and
// the correct response body, "OK".
func TestPing(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/healthcheck")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	t.Parallel()
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApp(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLS
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		// rebind tt into this lexical scope to avoid concurrency bug from running
		// sub-tests
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			code, _, body := ts.get(t, tt.urlPath)
			t.Logf("testing %q for want-code %d and want-body %q", tt.name, tt.wantCode,
				tt.wantBody)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q, but got %q", tt.wantBody, body)
			}
		})
	}

}

func TestSignupUser(t *testing.T) {
	// Create the application struct containing our mocked dependencies and
	// set up the test server for running an end-to-test.
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	fmt.Println(body)
	fmt.Printf("typ of body is %T", body)
	csrfToken := extractCSRFToken(t, body)

	// Log the CSRF token value in our test output. To see the output from the
	// t.Log() command you need to run `go test` with the -v (verbose) flag enabled.
	t.Log(csrfToken)
}
