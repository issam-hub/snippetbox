package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fs := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.AuthPassage)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.signUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.signUpPost))
	router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.RequireAuth)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.RecoverPanics, app.LogRequest, SecureHeaders)
	return standard.Then(router)
}
