package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a
		// panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			// If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				// This acts as a trigger to make Go's http server automatically close
				// the current connection after a response has been sent. It also informs
				// the user that the connection will be closed.
				w.Header().Set("Connection", "close")

				// Call the app.serverError helper method to return a 500 status code.
				// Also, panic returns an interface{}, so we normalize error into an
				// error object by using the fmt.Errorf() function (which app.serverError expects).
				// using fmt.Errorf() with err will create a new error object containing the default
				// textual representation of the interface{} value panic returns.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent hanlders in the chain are executed
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages that
		// require authentication are not stored in the users browsers cache (or other
		// intermediary cache)
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
