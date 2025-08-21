package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *application) updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data instead of JSON
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidInput)
		return
	}

	// Get username from form
	username := strings.TrimSpace(r.FormValue("username"))

	// Validate username
	if len(username) > 16 || len(username) == 0 {
		app.badRequestResponse(w, r, fmt.Errorf("username length must be between 1-16 characters"))
		return
	}

	user := app.getUserFromRequst(r)
	var imageURL string

	// Check if an image file was uploaded
	_, fileHeader, err := r.FormFile("image")
	if err == nil && fileHeader != nil {
		// Image was uploaded, process it
		imageURL, err = app.saveProfilePictureFromRequest(r)
		if err != nil {
			// Handle specific errors appropriately
			switch err {
			case ErrInvalidInput:
				app.badRequestResponse(w, r, fmt.Errorf("invalid image file"))
			case ErrFileSizeTooBig:
				app.badRequestResponse(w, r, fmt.Errorf("image file too large"))
			default:
				app.serverErrorResponse(w, r, err, "save profile picture")
			}
			return
		}

		// If user had an old profile picture, delete it
		if user.Image != "" {
			// Extract filename from old image URL
			oldFilename := strings.TrimPrefix(user.Image, "/media/profile-pictures/")
			if oldFilename != user.Image { // Make sure it's a valid profile picture URL
				_ = app.deleteProfilePicture(oldFilename) // Ignore error if file doesn't exist
			}
		}
	} else {
		// No new image uploaded, keep the existing one
		imageURL = user.Image
	}

	// Update user information in database
	err = app.userModel.UpdateImageAndUsername(user.ID, imageURL, username)
	if err != nil {
		app.serverErrorResponse(w, r, err, "update user info")
		return
	}

	app.writeJSON(w, 200, envelope{"message": "Updated successfully", "image_url": imageURL}, nil)
}
