package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"globechat.live/internal/models"
	"globechat.live/internal/random"
)

func (app *application) generateAccountObject(user models.User) map[string]any {
	return map[string]any{
		"id":          user.ID,
		"email":       user.Email,
		"username":    user.Username,
		"new_account": user.CreatedAt.Unix()-time.Now().Unix() < 10,
		"created_at":  user.CreatedAt.UTC(),
		"image":       user.Image,
		"messages":    user.Messages,
	}
}

func (app *application) createNewUser(email string) (models.User, error) {
	username := random.GenerateRandomUserName()
	// Username is available
	user, err := app.userModel.Create(email, username)
	return user, err
}

func (app *application) loginWihGoogleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSONFromRequest(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("nice try"))
		return
	}

	if strings.TrimSpace(input.Token) == "" {
		app.badRequestResponse(w, r, fmt.Errorf("are you trying to login without jwt token? are you fr?"))
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
		app.badRequestResponse(w, r, fmt.Errorf("nice try dude"))
		return
	}
	err = app.readJSON(res.Body, &responseData)

	if err != nil || res.StatusCode != 200 {
		app.serverErrorResponse(w, r, err, "read json")
		return
	}

	if responseData.Aud != app.config.googleClientId {
		app.badRequestResponse(w, r, fmt.Errorf("do you think you are smarter than me?"))
		return
	}

	if !slices.Contains([]string{"accounts.google.com", "https://accounts.google.com"}, responseData.Iss) {
		app.badRequestResponse(w, r, fmt.Errorf("is google drunk or are you doing something fishy?"))
		return
	}

	user, err := app.userModel.GetByEmail(responseData.Email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			user, err := app.createNewUser(responseData.Email)

			if err != nil {
				app.serverErrorResponse(w, r, err, "create user")
				return
			}
			token, err := app.sessionModel.Create(user.ID)

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

	token, err := app.sessionModel.Create(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err, "create session")
		return
	}
	app.writeJSON(w, 200, envelope{"token": token, "account": app.generateAccountObject(user)}, nil)
}

func (app *application) getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	app.writeJSON(w, 200, envelope{"account": app.generateAccountObject(*user)}, nil)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
			app.badRequestResponse(w, r, fmt.Errorf("you are trying to enter wrong terrority my guy"))
		}
	})
}
