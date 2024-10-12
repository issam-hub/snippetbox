package forms

import "snippetbox/interal/validator"

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}
