package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed html
var html embed.FS

func main() {
	// Create a file system from the embedded directory.
	staticFS, err := fs.Sub(html, "html")
	if err != nil {
		log.Fatal(err)
	}

	// Serve the files using the http.FileServer handler.
	http.Handle("/login/", http.StripPrefix("/login", http.FileServer(http.FS(staticFS))))

	// Start the HTTP server.
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
