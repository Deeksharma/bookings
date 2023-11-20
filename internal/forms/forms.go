package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form creates a custom form struct and, it embeds url.Values
type Form struct {
	url.Values
	Errors errors
}

// NewForm initialises a form struct
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required iterates through all the fields and add errors accordingly
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "*This field is required")
		}
	}
}

// Has checks is the post request has required parameters or not
func (f *Form) Has(field string) bool { // we'll have the http request and, it'll contain the all the id values
	fieldValue := f.Get(field) // useful in case of checkbox or if the field is present or not

	if fieldValue == "" {
		return false
	}
	return true
}

// Valid returns true if there is no errors else return false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks if the fields is long enough
func (f *Form) MinLength(field string, length int) bool { // we'll have the http request and, it'll contain the all the id values
	fieldValue := f.Get(field)
	if len(fieldValue) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail adds error if email is not valid
func (f *Form) IsEmail(field string)  {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, fmt.Sprintf("*Invalid email address"))
	}
}