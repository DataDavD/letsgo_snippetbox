package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/DataDavD/snippetbox/pkg/forms"
	"github.com/DataDavD/snippetbox/pkg/models"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		app.serverError(w, err)
		return
	}
}

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "home.page.gohtml", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a specific record
	// based on its ID. If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.gohtml", &templateData{
		Snippet: s,
	})
}

// createSnippetForm handler creates/renders snippet form response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.gohtml", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.NewForm(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any form_errors, we use our app.clientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the form,
	// then use the validation methods of forms.Form to check the content.
	form := forms.NewForm(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the form.Form object
	// as the data
	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded in the
	// form.Form struct, we can use the Get() method to retrieve the validated value for a
	// particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"),
		form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the session.Put() method to add a string value ("Your snippet was saved successfully")
	// and the corresponding key ("flash") to the session data. Note that if there is no existing
	// session for the current user (or their session has expired) then a new, empty,
	// session for them will automatically be created by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

	// Redirect the user to the relevant page for the created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

// signupUserForm handles the signup user form.
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", &templateData{
		Form: forms.NewForm(nil),
	})
}

// signupUser handles users getting signed up.
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using form helpsers.
	form := forms.NewForm(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password")) // Using embedded url.Values.Get method
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.FormErrors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.gohtml", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise, add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", &templateData{
		Form: forms.NewForm(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form errors map and re-display the login page.
	form := forms.NewForm(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"),
		form.Get("password")) // Using embedded url.Values.Get method
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.FormErrors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the ID of the current user to the session, so that they are now "logged in".
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
