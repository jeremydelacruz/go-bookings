package forms

import (
	"net/http"
	"net/url"
)

// Form defines a custom form struct
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and is not empty
func (f *Form) Has(field string, r *http.Request) bool {
	result := r.Form.Get(field)
	return result != ""
}
