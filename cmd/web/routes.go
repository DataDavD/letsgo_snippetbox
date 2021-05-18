package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard" middleware
	// which will be used for every request our app receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	// Register exact matches before wildcard route match (i.e. :id in Get method for
	// '/snippet/create')
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Wrap the existing chain with the recoverPanic middleware.
	return standardMiddleware.Then(mux)
}
