package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"environment": app.config.env,
		"status":      "I'm very okay 👌",
	}
	err := app.writeJSON(w, 200, envelope{"health": data}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err, "healthcheck json write")
	}
}
