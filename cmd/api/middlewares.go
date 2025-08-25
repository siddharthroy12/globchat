package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"globechat.live/internal/models"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err), "panic")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, exists := r.Header["Token"]

		if exists {
			if len(token) > 0 {
				user, err := a.userModel.GetFromSessionToken(token[0])
				if err == nil {
					ctx := context.WithValue(r.Context(), UserContextKey, &user)
					r = r.WithContext(ctx)
				} else {
					if !errors.Is(err, models.ErrNoRecord) {
						a.serverErrorResponse(w, r, err, "authenticate")
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
