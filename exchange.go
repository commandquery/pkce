package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
		Name:     config.Cookie,
		Value:    token.AccessToken,
		Expires:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		HttpOnly: true,
		Secure:   config.Secure,
		Domain:   config.Domain,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	fmt.Printf("setting cookie: %#v\n", cookie)

	// The URL that started the login process should be set in "state".
	redirectUrl := r.Form.Get("state")
	fmt.Println("redirecting state:", redirectUrl)
	if redirectUrl == "" {
		redirectUrl = config.Redirect
	}

	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
