package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/DataDavD/snippetbox/pkg/models"
)

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the curr req path exactly matches "/". If it doesn't, use
	// the http.NotFound() func to send 404 resp.
	if r.URL.Path != "/" {
		app.notFound(w)
		// return from func to avoid proceeding to home page response
		return
	}

	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	// Use the template.ParseFiles() func to read the files and store the templates
	// in a template set. Noticed that we can pass teh slice of file paths as a variadic
	// param?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract id value from query string and try to convert it to an integer
	// using strconv.Atoi() func. If it can't be converted, or the value is less than 1,
	// we return a 404 page Not Found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	// Write the snippet data as a plain-text HTTP response body.
	if _, err = fmt.Fprintf(w, "%v", s); err != nil {
		app.errorLog.Println("issue writing ID to response body", err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// only allow Post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on during the build
	title := "DataDavD Awesome Adventures in Life"
	content := "DataDavD has had an awesome, super, crazy, cool life!!!"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
