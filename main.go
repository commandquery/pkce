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
	Base     string // URL base used to access this service.
	Prefix   string // URL path prefix used to access this service.
	Issuer   string // The OIDC issuer URI
	ClientID string // Client ID as provided by the issuer
	Redirect string // application home page; redirect to this on success if no "state" given
	Domain   string // JWT Cookie domain
	Secure   bool   // Use secure cookies for JWT?
	Cookie   string // JWT Cookie name
	Log      bool   // Print some debug logging?
}

var client = &http.Client{}

var config = Config{
	Base:     "http://hello.example.com", // the PKCE service should probably live next to your auth services.
	Prefix:   "/pkce",                    // redirect_url will be BASE + PREFIX + "/authorize.html"
	ClientID: "",
	Issuer:   "https://example.com",
	Redirect: "http://app.example.com/home.", // The redirect should go a page in your app.
	Domain:   "example.com",
	Secure:   true,
	Cookie:   "Authorization",
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

func flag(def bool, env string) bool {
	value, ok := os.LookupEnv(env)
	if !ok {
		return def
	}

	return value == "true"
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

	config.Prefix = override(config.Prefix, "PKCE_PREFIX")
	config.ClientID = override(config.ClientID, "PKCE_CLIENT_ID")
	config.Base = override(config.Base, "PKCE_BASE")
	config.Issuer = override(config.Issuer, "PKCE_ISSUER")
	config.Domain = override(config.Domain, "PKCE_DOMAIN")
	config.Redirect = override(config.Redirect, "PKCE_REDIRECT")

	config.Secure = flag(true, "PKCE_SECURE")
	config.Log = flag(false, "PKCE_LOG")

	// Treat the root as if it's the login page, so that
	//     $PKCE_BASE/login will serve login.html
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
	if config.Log {
		fmt.Printf("listening on %s%s\n", addr, root)
	}

	// Start the HTTP server
	http.HandleFunc(root+"exchange", exchange)
	log.Fatal(http.ListenAndServe(addr, nil))
}
