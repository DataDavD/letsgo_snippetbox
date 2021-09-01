package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our app receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// add the authenticate() middleware to the chain
	// and use the noSurf middleware on all our dynamic routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()

	// Health check
	mux.Get("/healthcheck", http.HandlerFunc(app.ping))

	// Register exact matches before wildcard route match (i.e. :id in Get method for
	// '/snippet/create').
	// Update these routes to use the dynamic middleware chain follow by the appropriate handler
	// function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// Require auth middleware for auth'd/logged-in actions
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	// Require auth middleware for auth'd/logged-in actions
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippet))

	// Add the five new routes for user authentication.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	// Require auth middleware for auth'd/logged-in actions
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
