package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form anonymously embeds a url.Values object to (to hold the form data)
// and an FormErrors field to hold any validation errors for the form data
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

// Valid method checks FormErrors for any present errors. It returns true if there are no errors,
// else it returns false if there are errors.
func (f *Form) Valid() bool {
	return len(f.FormErrors) == 0
}
