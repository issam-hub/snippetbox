package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/interal/forms"
	"snippetbox/interal/models"
	templates "snippetbox/ui/static/templates"
	"strconv"

	"snippetbox/interal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	component := templates.Home(snippets)
	component.Render(r.Context(), w)

}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}
	component := templates.View(*snippet)
	component.Render(r.Context(), w)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	form := forms.SnippetCreateForm{}
	component := templates.CreateSnippet(form)
	component.Render(r.Context(), w)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.SnippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title cannot be more than 100 characters")

	form.CheckField(validator.NotBlank(form.Content), "content", "Content cannot be blank")

	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "Expires can be only one of the following values: 1, 7, 365")

	if !form.Valid() {
		component := templates.CreateSnippet(form)
		component.Render(r.Context(), w)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
