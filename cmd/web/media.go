package main

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func (app *application) saveProfilePicture(img image.Image, filename string) error {
	// Create profile pictures directory path
	profilePicturesDir := filepath.Join(app.config.mediaDir, "profile-pictures")

	// Ensure the profile pictures directory exists
	if err := os.MkdirAll(profilePicturesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile pictures directory: %w", err)
	}

	// Create the full file path
	filePath := filepath.Join(profilePicturesDir, filename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode and save as JPEG
	jpegOptions := &jpeg.Options{Quality: 90}
	if err := jpeg.Encode(file, img, jpegOptions); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

func (app *application) deleteProfilePicture(filename string) error {
	// Construct the full file path in profile pictures directory
	profilePicturesDir := filepath.Join(app.config.mediaDir, "profile-pictures")
	filePath := filepath.Join(profilePicturesDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (app *application) saveProfilePictureFromRequest(r *http.Request) (string, error) {
	// Parse multipart form with 32MB max memory
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", ErrInvalidInput
	}

	// Get the file from form data
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", ErrInvalidInput
	}
	defer file.Close()

	// Validate file size (optional - limit to 10MB)
	if header.Size > 10<<20 {
		return "", ErrFileSizeTooBig
	}

	// Decode the image (supports JPEG, PNG, GIF)
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Resize image to 38x38 pixels
	resizedImg := resize.Resize(76, 76, img, resize.Lanczos3)

	// Generate random filename
	randomID, err := generateRandomID(16) // 32 character hex string
	if err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s.jpg", randomID)

	// Save the resized image to file
	if err := app.saveProfilePicture(resizedImg, filename); err != nil {
		return "", err
	}

	// Construct the public URL
	// Assuming you have a config field for base URL and media endpoint
	// Adjust this according to your application's URL structure
	url := fmt.Sprintf("/media/profile-pictures/%s", filename)
	return url, nil
}

// mediaHandler serves static media files
func (app *application) mediaHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the file path from the URL
	// Remove "/media/" prefix to get the relative path
	filePath := strings.TrimPrefix(r.URL.Path, "/media/")

	// Prevent directory traversal attacks
	if strings.Contains(filePath, "..") {
		app.notFoundHandler(w, r)
		return
	}

	// Construct full file path
	fullPath := filepath.Join(app.config.mediaDir, filePath)

	// Check if file exists and is not a directory
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		app.notFoundHandler(w, r)
		return
	}

	if fileInfo.IsDir() {
		app.notFoundHandler(w, r)
		return
	}

	// Set appropriate content type based on file extension
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Set cache headers for better performance
	w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
	w.Header().Set("ETag", fmt.Sprintf(`"%x-%x"`, fileInfo.ModTime().Unix(), fileInfo.Size()))

	// Serve the file
	http.ServeFile(w, r, fullPath)
}
