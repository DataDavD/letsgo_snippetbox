package main

import (
	"log"
	"net/http"
)

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func home(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello from Snippetbox")); err != nil {
		log.Println("home request:", err)
	}
}

func main() {
	// Use the http.NewServeMux() func to tinit a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

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
