package main

import (
	"fmt"
	"net/http"
	"strings"

	"globechat.live/internal/models"
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

func (app *application) queryUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Get search parameter
	search := strings.TrimSpace(queryParams.Get("search"))

	// Get page size with default and validation
	pageSize, err := app.readInt(queryParams, "page_size", 20) // default 20
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid page_size parameter: %v", err))
		return
	}
	if pageSize < 1 || pageSize > 100 {
		app.badRequestResponse(w, r, fmt.Errorf("page_size must be between 1 and 100"))
		return
	}

	// Get page index with default and validation
	pageIndex, err := app.readInt(queryParams, "page", 0) // default 0 (first page)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid page parameter: %v", err))
		return
	}
	if pageIndex < 0 {
		app.badRequestResponse(w, r, fmt.Errorf("page must be 0 or greater"))
		return
	}

	// Create query struct
	query := models.UserQuery{
		Search:    search,
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}

	// Execute query
	result, err := app.userModel.Query(query)
	if err != nil {
		app.serverErrorResponse(w, r, err, "query users")
		return
	}

	// Calculate pagination metadata
	totalPages := 0
	if pageSize > 0 {
		totalPages = (result.Total + pageSize - 1) / pageSize // ceiling division
	}

	hasNext := pageIndex < totalPages-1
	hasPrev := pageIndex > 0

	// Prepare response
	response := envelope{
		"users": result.Users, // Note: assuming UserQueryResult was fixed to have Users field
		"pagination": envelope{
			"total":       result.Total,
			"count":       result.Count,
			"page":        pageIndex,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	}

	// Add search info if search was performed
	if search != "" {
		response["search"] = search
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}
