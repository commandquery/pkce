package main

import (
	"embed"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// Config can be overridden by environment variables
type Config struct {
	Prefix   string
	ClientID string
	Base     string
	Issuer   string
	Redirect string
	Domain   string
	Secure   bool
}

var client = &http.Client{}

var config = Config{
	Prefix:   "",
	ClientID: "205879184755607046@bookwork",
	Base:     "http://localhost:8081",
	Issuer:   "https://hello.bookwork.com",
	Redirect: "http://localhost:8081/home.html",
	Domain:   "localhost",
	Secure:   false,
}

//go:embed templates
var embedFS embed.FS
var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(embedFS, "templates/*")
	if err != nil {
		log.Fatal(err)
	}
}

// Override a value with an environment variable, if it's defined.
func override(def string, env string) string {
	value := os.Getenv(env)
	if value != "" {
		return value
	} else {
		return def
	}
}

func writeTemplate(w http.ResponseWriter, path string) {
	tmpl := templates.Lookup(path)
	if tmpl == nil {
		http.Error(w, path, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))

	// Execute the template with the input data
	if err := tmpl.Execute(w, &config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	config.Prefix = override(config.Prefix, "JWT_PREFIX")
	config.ClientID = override(config.ClientID, "JWT_CLIENT_ID")
	config.Base = override(config.Base, "JWT_BASE")
	config.Issuer = override(config.Issuer, "JWT_ISSUER")
	config.Domain = override(config.Domain, "JWT_DOMAIN")
	config.Redirect = override(config.Redirect, "JWT_REDIRECT")

	if os.Getenv("JWT_SECURE") == "true" {
		config.Secure = true
	}

	// Treat the root as if it's the login page, so that
	//     hello.bookwork.com/login will serve login.html
	// All other pages need to be named specifically.
	// This simply makes it easier to redirect to a login page.
	root := config.Prefix + "/"
	http.HandleFunc(root, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == root {
			writeTemplate(w, "login.html")
		} else {
			writeTemplate(w, filepath.Base(r.URL.Path))
		}
	})

	addr := ":8081"
	fmt.Printf("listening on %s%s\n", addr, root)
	// Start the HTTP server
	http.HandleFunc(root+"exchange", exchange)
	log.Fatal(http.ListenAndServe(addr, nil))
}
