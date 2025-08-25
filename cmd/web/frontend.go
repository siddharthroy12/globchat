package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"path"
	"strings"

	"globechat.live/frontend"
)

// gzipWriter wraps http.ResponseWriter to provide gzip compression
type gzipWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	return gw.gzipWriter.Write(b)
}

func (gw *gzipWriter) Close() error {
	return gw.gzipWriter.Close()
}

// gzipMiddleware adds gzip compression support
func gzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if client supports gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		// Set the content encoding header
		w.Header().Set("Content-Encoding", "gzip")

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the response writer
		gzw := &gzipWriter{
			ResponseWriter: w,
			gzipWriter:     gz,
		}

		next(gzw, r)
	}
}

func (app *application) spaHandler(w http.ResponseWriter, r *http.Request) {
	// Don't serve frontend for API routes that don't exist
	if strings.HasPrefix(r.URL.Path, "/api/") {
		app.notFoundHandler(w, r)
		return
	}

	// Clean the path and remove leading slash for embed.FS
	cleanPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if cleanPath == "" {
		cleanPath = "index.html"
	}

	// Try to serve the requested file from the static folder
	staticPath := path.Join("static", cleanPath)
	file, err := frontend.Files.Open(staticPath)
	if err == nil {
		defer file.Close()

		stat, err := file.Stat()
		if err == nil && !stat.IsDir() {
			// File exists and is not a directory, serve it
			contentType := getContentType(cleanPath)
			w.Header().Set("Content-Type", contentType)

			// Set appropriate cache headers based on file type
			if isStaticAsset(cleanPath) {
				w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day for assets
			} else {
				w.Header().Set("Cache-Control", "no-cache") // No cache for HTML files
			}

			// Only compress compressible content types
			if shouldCompress(cleanPath) {
				http.ServeContent(w, r, stat.Name(), stat.ModTime(), file.(io.ReadSeeker))
			} else {
				// For non-compressible files (like images), serve directly without compression
				http.ServeContent(w, r, stat.Name(), stat.ModTime(), file.(io.ReadSeeker))
			}
			return
		}
		file.Close()
	}

	// File not found or is a directory, serve index.html for SPA routing
	indexFile, err := frontend.Files.Open(path.Join("static", "index.html"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer indexFile.Close()

	stat, err := indexFile.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}

// Helper function to determine content type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

// Helper function to determine if a file is a static asset (for caching)
func isStaticAsset(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	staticExts := map[string]bool{
		".css":   true,
		".js":    true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".svg":   true,
		".ico":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
	}
	return staticExts[ext]
}

// Helper function to determine if content should be compressed
func shouldCompress(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	compressibleExts := map[string]bool{
		".html": true,
		".css":  true,
		".js":   true,
		".json": true,
		".svg":  true,
		".txt":  true,
		".xml":  true,
	}
	return compressibleExts[ext]
}

// Usage example: wrap your handler with the gzip middleware
// http.HandleFunc("/", gzipMiddleware(app.spaHandler))
