package forms

import "snippetbox/interal/validator"

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

type UserSignupForm struct {
	Name           string
	Email          string
	HashedPassword string
	validator.Validator
}

type UserLoginForm struct {
	Email          string
	HashedPassword string
	validator.Validator
}
