package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"globechat.live/internal/models"
)

func (app *application) generateAccountObject(user models.User) map[string]string {
	return map[string]string{
		"email":    user.Email,
		"username": user.Username,
	}
}

func (app *application) loginWihGoogle(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSONFromRequest(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if strings.TrimSpace(input.Token) == "" {
		app.badRequestResponse(w, r, fmt.Errorf("token is empty"))
		return
	}

	var responseData struct {
		Aud   string `json:"aud"`
		Iss   string `json:"iss"`
		Email string `json:"email"`
	}

	res, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", input.Token))

	if err != nil {
		app.serverErrorResponse(w, r, err, "google login api call")
		return
	}

	if res.StatusCode != 200 {
		app.badRequestResponse(w, r, ErrInvalidToken)
		return
	}
	err = app.readJSON(res.Body, &responseData)

	if err != nil || res.StatusCode != 200 {
		app.serverErrorResponse(w, r, err, "read json")
		return
	}

	if responseData.Aud != app.config.googleClientId {
		app.badRequestResponse(w, r, ErrInvalidToken)
		return
	}

	if !slices.Contains([]string{"accounts.google.com", "https://accounts.google.com"}, responseData.Iss) {
		app.badRequestResponse(w, r, ErrInvalidToken)
		return
	}

	user, err := app.userModel.GetByEmail(responseData.Email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// Create account

			userId, err := app.userModel.Create(responseData.Email, responseData.Email)

			if err != nil {
				app.serverErrorResponse(w, r, err, "create user")
				return
			}
			token, err := app.sessionModel.Create(userId)

			if err != nil {
				app.serverErrorResponse(w, r, err, "create session")
				return
			}
			app.writeJSON(w, 200, envelope{"token": token, "account": app.generateAccountObject(user)}, nil)
			return
		} else {
			app.serverErrorResponse(w, r, err, "get user by email")
			return
		}
	}

	// Create another session and return it
	token, err := app.sessionModel.Create(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err, "create session")
		return
	}
	app.writeJSON(w, 200, envelope{"token": token, "account": app.generateAccountObject(user)}, nil)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	err := app.sessionModel.Remove(user.ID)

	if err != nil {
		app.serverErrorResponse(w, r, err, "session remove")
		return
	}
	app.writeJSON(w, 200, envelope{"message": "Successfully logged out"}, nil)
}

func (app *application) getUserFromRequst(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		panic(fmt.Errorf("trying to access user for path for which authentication is not required"))
	}

	return user
}

func (app *application) isAuthenticated(r *http.Request) bool {
	_, ok := r.Context().Value(UserContextKey).(*models.User)

	return ok
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			next(w, r)
			return
		} else {
			app.badRequestResponse(w, r, ErrInvalidToken)
		}
	})
}
