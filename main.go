package main

import (
	"embed"
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
	Base:     "http://localhost:8080",
	Issuer:   "https://hello.bookwork.com",
	Redirect: "http://localhost:8080/home.html",
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
func override(env string, def string) string {
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

	config.Prefix = override("JWT_PREFIX", config.Prefix)
	config.ClientID = override("JWT_CLIENT_ID", config.ClientID)
	config.Base = override("JWT_BASE", config.Base)
	config.Issuer = override("JWT_ISSUER", config.Issuer)
	config.Domain = override("JWT_DOMAIN", config.Domain)
	config.Redirect = override("JWT_REDIRECT", config.Redirect)

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

	// Start the HTTP server
	http.HandleFunc(config.Prefix+"/exchange", exchange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
