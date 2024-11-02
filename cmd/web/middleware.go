package main

import (
	"context"
	"fmt"
	"net/http"
	"snippetbox/interal/models"

	"github.com/justinas/nosurf"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, req)
	})
}

func (app *application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())

		next.ServeHTTP(w, req)
	})
}

func (app *application) RecoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, req)
	})
}

func (app *application) AddTemplateData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		crosstemplates := &models.CrossTemplates{
			IsAuthenticated: app.isAuthenticated(req),
			CSRFToken:       nosurf.Token(req),
		}
		ctx := context.WithValue(req.Context(), models.ContextClass, crosstemplates)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func (app *application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !app.isAuthenticated(req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, req)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		// Secure:   true,
	})

	return csrfHandler
}
