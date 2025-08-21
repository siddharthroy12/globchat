package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// API routes
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/google/login", app.loginWihGoogleHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/user", app.requireAuthentication(app.getUserDataHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/user", app.requireAuthentication(app.updateUserInfoHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/logout", app.requireAuthentication(app.logoutHandler))

	// Media files route (serves static files)
	router.HandlerFunc(http.MethodGet, "/media/*filepath", app.mediaHandler)

	// Error handlers
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedHandler)

	return app.recoverPanic(app.authenticate(router))
}
