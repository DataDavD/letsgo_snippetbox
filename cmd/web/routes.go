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

	// Create a new middleware chain containing the specific to our dynamic application routes.
	// For now, this chain will only contain the session middleware, but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	// Register exact matches before wildcard route match (i.e. :id in Get method for
	// '/snippet/create').
	// Update these routes to use the dynamic middleware chain follow by the appropriate handler
	// function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
