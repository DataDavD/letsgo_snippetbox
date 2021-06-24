package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile() function to parse a pattern and compile a
// regular expression for sanity checking the format of an email address.
// This returns a *regexp.Regexp object, or panics in the event of an error.
// Doing this once at runtime, and storing the compiled regular expression
// object in a variable, is more performant than re-compiling the pattern with
// every request.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](" +
	"?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form anonymously embeds a url.Values object
// to (to hold the form data) and an FormErrors field to hold any validation errors
// for the form data
type Form struct {
	url.Values
	FormErrors formErrors
}

// NewForm initializes a custom Form struct. Notice that this takes the form data as
// the parameter.
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		formErrors(map[string][]string{}),
	}
}

// Required methods checks that specific fields in the form data are present and not blank
// If any fields fail this check, add the appropriate message to the FormErrors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field) // Using embedded url.Values.Get method
		if strings.TrimSpace(value) == "" {
			f.FormErrors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks that a specific field in the form contains a maximum number of characters.
// If the check fails then it adds the appropriate message to FormErrors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field) // Using embedded url.Values.Get method
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.FormErrors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters", d))
	}
}

// PermittedValues method checks that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message the FormErrors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field) // Using embedded url.Values.Get method
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.FormErrors.Add(field, "This field is invalid")
}

// MinLength checks that a specific field in the form contains a minimum number of characters.
// If the check fails it adds the appropriate message to the form errors.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.FormErrors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)",
			d))
	}
}

// MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.FormErrors.Add(field, "This field is invalid")
	}
}

// Valid method checks FormErrors for any present errors. It returns true if there are no errors,
// else it returns false if there are errors.
func (f *Form) Valid() bool {
	return len(f.FormErrors) == 0
}
