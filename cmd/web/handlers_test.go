package main

import (
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
