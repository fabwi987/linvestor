package models

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	clientID        = "671169303456-8qk6ut70hnpu3iig62soqmj53nhs1daa.apps.googleusercontent.com"
	clientSecret    = "cxJf-iZ4G9n2Q1PUedajavS-"
	applicationName = "Linvestor"
)

var config = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  "postmessage",
	Scopes: []string{
		"https://www.googleapis.com/auth/plus.login", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	},
	Endpoint: google.Endpoint,
}

var store = sessions.NewCookieStore([]byte(randomString(32)))

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

type ClaimSet struct {
	Sub string
}

func exchange(code string) (accessToken string, idToken string, err error) {
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", "", fmt.Errorf("Error while exchanging code: %v", err)
	}
	// TODO: return ID token in second parameter from updated oauth2 interface
	return tok.AccessToken, tok.Extra("id_token").(string), nil
}

func decodeIdToken(idToken string) (gplusID string, err error) {
	// An ID token is a cryptographically-signed JSON object encoded in base 64.
	// Normally, it is critical that you validate an ID token before you use it,
	// but since you are communicating directly with Google over an
	// intermediary-free HTTPS channel and using your Client Secret to
	// authenticate yourself to Google, you can be confident that the token you
	// receive really comes from Google and is valid. If your server passes the ID
	// token to other components of your app, it is extremely important that the
	// other components validate the token before using it.
	var set ClaimSet
	if idToken != "" {
		// Check that the padding is correct for a base64decode
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Malformed ID token")
		}
		// Decode the ID token
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
	}
	return set.Sub, nil
}

// AuthenticationHandler sets up a session for the current user and serves the home page
func AuthenticationHandler(c *gin.Context) {
	// This check prevents the "/" handler from handling all requests by default

	/**if c.Request.URL.Path != "/" {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Invalid session state: %s", state))
		return
	}*/

	// Create a state token to prevent request forgery and store it in the session
	// for later validation
	session, err := store.Get(c.Request, "sessionName")
	log.Println("Session: [" + session.Name() + "]")
	if err != nil {
		log.Println("error fetching session:", err)
		// Ignore the initial session fetch error, as Get() always returns a
		// session, even if empty.
		//return &appError{err, "Error fetching session", 500}
	}

	state := randomString(64)
	session.Values["state"] = state
	session.Save(c.Request, c.Writer)

	stateURL := url.QueryEscape(session.Values["state"].(string))
	log.Println("State URL: [" + stateURL + "]")

	// Fill in the missing fields in index.html
	var data = struct {
		ApplicationName, ClientID, State string
	}{applicationName, clientID, stateURL}

	log.Println("Data: [" + data.ApplicationName + "]")
	log.Println("Data: [" + data.ClientID + "]")
	log.Println("Data: [" + data.State + "]")
	// Render and serve the HTML

	c.HTML(http.StatusOK, "google.html", data)

}

// connect exchanges the one-time authorization code for a token and stores the
// token in the session
func Connect(c *gin.Context) {
	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)

	}
	if c.PostForm("state") != session.Values["state"].(string) {
		log.Println("Invalid state parameter")
	}
	// Normally, the state is a one-time token; however, in this example, we want
	// the user to be able to connect and disconnect without reloading the page.
	// Thus, for demonstration, we don't implement this best practice.
	// session.Values["state"] = nil

	// Setup for fetching the code from the request payload
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//return &appError{err, "Error reading code in request body", 500}
	}
	code := string(x)

	accessToken, idToken, err := exchange(code)
	if err != nil {
		//return &appError{err, "Error exchanging code for access token", 500}
	}
	gplusID, err := decodeIdToken(idToken)
	if err != nil {
		//	return &appError{err, "Error decoding ID token", 500}
	}

	// Check if the user is already connected
	storedToken := session.Values["accessToken"]
	storedGPlusID := session.Values["gplusID"]
	if storedToken != nil && storedGPlusID == gplusID {
		//m := "Current user already connected"
		log.Println("Current user already connected")
		//return &appError{errors.New(m), m, 200}
	}

	// Store the access token in the session for later use
	session.Values["accessToken"] = accessToken
	session.Values["gplusID"] = gplusID
	session.Save(c.Request, c.Writer)
	//return nil
}

// disconnect revokes the current user's token and resets their session
func Disconnect(c *gin.Context) {
	// Only disconnect a connected user
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		//return &appError{err, "Error fetching session", 500}
	}
	token := session.Values["accessToken"]
	if token == nil {
		//m := "Current user not connected"
		log.Println("Current user not connected", err)
		//return &appError{errors.New(m), m, 401}
	}

	// Execute HTTP GET request to revoke current token
	url := "https://accounts.google.com/o/oauth2/revoke?token=" + token.(string)
	resp, err := http.Get(url)
	if err != nil {
		//m := "Failed to revoke token for a given user"
		log.Println("Failed to revoke token for a given user", err)
		//return &appError{errors.New(m), m, 400}
	}
	defer resp.Body.Close()

	// Reset the user's session
	session.Values["accessToken"] = nil
	session.Save(c.Request, c.Writer)
	//return nil
}

func Loggedin(c *gin.Context) {
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		//return &appError{err, "Error fetching session", 500}
	}
	token := session.Values["accessToken"]
	// Only fetch a list of people for connected users
	if token == nil {
		//m := "Current user not connected"
		log.Println("Current user not connected", err)
		//return &appError{errors.New(m), m, 401}
	}

	/**
	// Create a new authorized API client
	tok := new(oauth2.Token)
	tok.AccessToken = token.(string)
	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(tok))
	service, err := plus.New(client)
	if err != nil {
		log.Println("Plus error ", err)
		//return &appError{err, "Create Plus Client", 500}
	}

	// Get a list of people that this user has shared with this app

	people := service.People.List("me", "visible")
	peopleFeed, err := people.Do()
	if err != nil {
		m := "Failed to refresh access token"
		if err.Error() == "AccessTokenRefreshError" {
			return &appError{errors.New(m), m, 500}
		}
		return &appError{err, m, 500}
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(&peopleFeed)
	if err != nil {
		return &appError{err, "Convert PeopleFeed to JSON", 500}
	}
	return nil*/
}

// randomString returns a random string with the specified length
func randomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
