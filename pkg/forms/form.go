package forms

import (
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	Entries url.Values
	// Slice of strings to enable multiple flash error messages for the same field
	Errors map[string][]string
}

func NewForm(r *http.Request) *Form {
	form := &Form{
		Entries: r.PostForm,
		Errors: make(map[string][]string),
	}

	return form
}

// Append a message to the slice of error messages for a given field
func (f *Form) AppendError(field, message string) {
	f.Errors[field] = append(f.Errors[field], message)
}

// Return only the first error string for a given field, if one exists
func (f *Form) GetFirstError(field string) string {
	messages := f.Errors[field]
	if len(messages) == 0 {
		return ""
	}
	return messages[0]
}

// Return all errors for a given field, if any exist
func (f *Form) GetAllErrors(field string) []string {
	messages := f.Errors[field]

	return messages
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Entries.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors[field] = append(f.Errors[field], value)
		}
	}
}