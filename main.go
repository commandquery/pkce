package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed templates
var embedFS embed.FS

// Config can be overridden by environment variables
type Config struct {
	ClientID string
	Base     string
	Issuer   string
	Redirect string
	Domain   string
	Secure   bool
}

var client = &http.Client{}

// Config for production.
// This is stored as overriding environment variables in deploy.yaml.
//var prodConfig = Config{
//	ClientID: "205769478103975430@bookwork",
//	Base:     "https://hello.bookwork.com/login",
//	Issuer:   "https://hello.bookwork.com",
//}

var config = Config{
	ClientID: "205879184755607046@bookwork",
	Base:     "http://localhost:8080",
	Issuer:   "https://hello.bookwork.com",
	Redirect: "http://localhost:8080/home.html",
	Domain:   "localhost",
	Secure:   false,
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

// Accept the login code and other parameters from the client, as well as the
// challenge string, and call Zitadel to complete the PKCE login process.
// We then set a HTTPOnly cookie containing the JWT.
//
// This needs to be done on the server side, because we can't set a HTTPOnly
// cookie on the client side.
func exchange(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", r.Form.Get("code"))
	data.Set("redirect_uri", config.Base+"/authorize.html") // this is the URI that received the redirect. security feature?
	data.Set("client_id", config.ClientID)
	data.Set("code_verifier", r.Form.Get("code_verifier"))

	// Encode the form data
	body := strings.NewReader(data.Encode())

	// Create the request
	req, err := http.NewRequest("POST", config.Issuer+"/oauth/v2/token", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create the HTTP client and send the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		IdToken     string `json:"id_token"`
	}{}

	err = json.Unmarshal(b, &token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token.AccessToken,
		Expires:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		HttpOnly: true,
		Secure:   config.Secure,
		Domain:   config.Domain,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, config.Redirect, http.StatusFound)
}

func main() {
	templates, err := fs.Sub(embedFS, "templates")
	if err != nil {
		log.Fatal(err)
	}

	config.ClientID = override("JWT_CLIENT_ID", config.ClientID)
	config.Base = override("JWT_BASE", config.Base)
	config.Issuer = override("JWT_ISSUER", config.Issuer)
	config.Domain = override("JWT_DOMAIN", config.Domain)
	config.Redirect = override("JWT_REDIRECT", config.Redirect)

	if os.Getenv("JWT_SECURE") == "true" {
		config.Secure = true
	}

	// Define the HTTP handler that will use the templates
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the template with the given name
		tmpl, err := template.ParseFS(templates, r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))

		// Execute the template with the input data
		if err := tmpl.Execute(w, &config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	http.HandleFunc("/exchange", exchange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
