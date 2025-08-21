package main

import (
	"fmt"
	"net/http"
)

var ErrInvalidToken = fmt.Errorf("invalid token")
var ErrFileSizeTooBig = fmt.Errorf("file size is too big")
var ErrInvalidInput = fmt.Errorf("invalid input")

func (app *application) logError(r *http.Request, err error, action string) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "action", action)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err, "errorResponse")
		w.WriteHeader(500)
	}
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error, action string) {
	app.logError(r, err, action)

	message := "this server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := "the requested resouces could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusNotFound, message)
}
