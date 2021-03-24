package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the curr req path exactly matches "/". If it doesn't, use
	// the http.NotFound() func to send 404 resp.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		// return from func to avoid proceeding to home page response
		return
	}
	if _, err := w.Write([]byte("Hello from Snippetbox")); err != nil {
		log.Println("home request:", err)
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract id value from query string and try to convert it to an integer
	// using strconv.Atoi() func. If it can't be converted, or the value is less than 1,
	// we return a 404 page Not Found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	if _, err := fmt.Fprintf(w, "Display a specific snippet with ID %d", id); err != nil {
		log.Println("show snippet request:", err)
	}
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	// only allow Post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", 405)
		return
	}

	if _, err := w.Write([]byte("Create a new snippet")); err != nil {
		log.Println("create new snippet request:", err)
	}
}

func main() {
	// Use the http.NewServeMux() func to tinit a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.

	log.Println("Starting server on localhost:4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
