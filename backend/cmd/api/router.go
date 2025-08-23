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
	router.HandlerFunc(http.MethodGet, "/api/v1/threads", app.getThreadsHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/threads/random", app.getRandomThread)
	router.HandlerFunc(http.MethodPost, "/api/v1/threads", app.requireAuthentication(app.createThreadHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/threads", app.requireAuthentication(app.deleteThreadHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/messages", app.requireAuthentication(app.createMessageHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/messages", app.requireAuthentication(app.deleteMessageHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/messages", app.getMessagesHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/messages/report", app.requireAuthentication(app.reportMessageHandler))

	// Websocket
	router.HandlerFunc(http.MethodGet, "/api/v1/ws", app.websocketConnectionHandler)

	// Media files route (serves static files)
	router.HandlerFunc(http.MethodGet, "/media/*filepath", app.mediaHandler)

	// Error handlers
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedHandler)

	return app.recoverPanic(app.authenticate(router))
}
