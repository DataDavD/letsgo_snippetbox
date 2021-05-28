package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/DataDavD/snippetbox/pkg/models"
)

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
	app.render(w, r, "show.page.gohtml", &templateData{Snippet: s})
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.gohtml", nil)
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

	// Use the r.PostGorm.Get() method to retrieve the relevant data fields from the
	// the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Initialize a map to hold any validation form_errors.
	formErrors := make(map[string]string)

	// Check that the title field is not blank and is not more than 100 characters
	// long. If it fails either of those checks, add a message to the form_errors
	// map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		formErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		formErrors["title"] = "This field is too long (maximum is 100 characters"
	}

	// check that the Content field isn't blank
	if strings.TrimSpace(content) == "" {
		formErrors["content"] = "This field cannot be blank"
	}

	// Check that the expires field isn't blank and matches one of the permitted
	// values ("1", "7", "365").
	if strings.TrimSpace(expires) == "" {
		formErrors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		formErrors["expires"] = "This field is invalid"
	}

	// If there are any form_errors, dump them in a plain text HTTP response and return
	// from handler.
	if len(formErrors) > 0 {
		if _, err := fmt.Fprint(w, formErrors); err != nil {
			app.serverError(w, err)
			return
		}
		return
	}

	// Create a new snippet record in the database using the form data.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
