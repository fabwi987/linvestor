package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//Google API project information
const (
	clientID        = "671169303456-8qk6ut70hnpu3iig62soqmj53nhs1daa.apps.googleusercontent.com"
	clientSecret    = "cxJf-iZ4G9n2Q1PUedajavS-"
	applicationName = "Google+ Go Quickstart"
)

// config is the configuration specification supplied to the OAuth package.
var config = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	// Scope determines which API calls you are authorized to make
	Scopes:   []string{"https://www.googleapis.com/auth/plus.login"},
	Endpoint: google.Endpoint,
	// Use "postmessage" for the code-flow for server side apps
	RedirectURL: "postmessage",
}

// store initializes the Gorilla session store.
var store = sessions.NewCookieStore([]byte(randomString(32)))

// Token represents an OAuth token response.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

// ClaimSet represents an IdToken response.
type ClaimSet struct {
	Sub string
}

// exchange takes an authentication code and exchanges it with the OAuth
// endpoint for a Google API bearer token and a Google+ ID
func exchange(code string) (accessToken string, idToken string, err error) {
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", "", fmt.Errorf("Error while exchanging code: %v", err)
	}
	// TODO: return ID token in second parameter from updated oauth2 interface
	return tok.AccessToken, tok.Extra("id_token").(string), nil
}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within
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

//InitiateSession creates a session and serves the page to the client
func InitiateSession(c *gin.Context) {
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		// Ignore the initial session fetch error, as Get() always returns a
		// session, even if empty.
		//return &appError{err, "Error fetching session", 500}
	}

	state := randomString(64)
	session.Values["state"] = state
	stateURL := url.QueryEscape(session.Values["state"].(string))
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println("Something went wrong when saving the session: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// Fill in the missing fields in index.html

	var data = struct {
		ApplicationName, ClientID, State string
	}{applicationName, clientID, stateURL}

	c.HTML(http.StatusOK, "google.html", data)
	//c.String(http.StatusOK, "Session state: %s Session stateUrl: %s", state, stateURL)

}

//Connect handles the exchange with Google
func Connect(c *gin.Context) {
	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	log.Println("Value in post form:")
	log.Println(c.Query("state"))
	if c.Query("state") != session.Values["state"].(string) {
		log.Println("Invalid state parameter")
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	//TODO: make sure that the code is sent from the client
	// Get the access code from Googel via the client
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	code := string(x)
	log.Println("Code: " + code)

	//Perform the exchange of the accessToken with Google
	accessToken, idToken, err := exchange(code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	//Decode the Google Plus ID
	gplusID, err := decodeIdToken(idToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	//Save the token and the ID in our session
	session.Values["accessToken"] = accessToken
	session.Values["gplusID"] = gplusID
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println("Something went wrong when saving the session: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	/**
	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	dbtable := "stock125"
	stmt, err := db.Prepare("INSERT " + dbtable + "SET users=?,googleid=?,name=?")
	Perror(err)

	res, err := stmt.Exec(1, gplusID, "Fabian Widen")
	Perror(err)
	nid, err := res.LastInsertId()
	log.Println(nid)

	db.Close()*/

}

//Disconnect removes the current users data from the session
func Disconnect(c *gin.Context) {
	// Only disconnect a connected user
	session, err := store.Get(c.Request, "sessionName")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	token := session.Values["accessToken"]
	log.Println(token)
	if token == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Execute HTTP GET request to revoke current token
	url := "https://accounts.google.com/o/oauth2/revoke?token=" + token.(string)
	resp, err := http.Get(url)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	defer resp.Body.Close()

	//Reset the session values
	session.Values["accessToken"] = nil
	session.Save(c.Request, c.Writer)

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
