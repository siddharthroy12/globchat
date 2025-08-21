package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *application) updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Username string `json:"username"`
		Image    string `json:"image"`
	}{}

	err := app.readJSONFromRequest(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.Username = strings.TrimSpace(input.Username)

	if len(input.Username) > 16 || len(input.Username) == 0 {
		app.badRequestResponse(w, r, fmt.Errorf("username length must be between 0 - 16"))
		return
	}
	user := app.getUserFromRequst(r)

	app.userModel.UpdateImageAndUsername(user.ID, input.Image, input.Username)

	app.writeJSON(w, 200, envelope{"message": "Updated successfully"}, nil)
}
