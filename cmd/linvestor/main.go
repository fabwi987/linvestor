package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/fabwi987/linvestor/controllers"
	"github.com/fabwi987/linvestor/models"
	"github.com/gin-gonic/gin"
)

var checker bool

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" //GG
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	//testetst
	router.GET("/start", ginFunc)
	router.GET("/old", oldStockView)
	router.GET("/new", insertStock)
	router.POST("/submit", submit)
	router.POST("/sell", sell)
	router.Run(":" + port)

}

//ginFunc returns a gin context
func ginFunc(c *gin.Context) {

	var dispData []models.StockDataSaveFormat
	dispData, growth := controllers.ShowStocks("stock123")

	layoutData := struct {
		ThreadID int
		Posts    []models.StockDataSaveFormat
		Growth   string
	}{
		ThreadID: 1,
		Posts:    dispData,
		Growth:   growth,
	}

	c.HTML(http.StatusOK, "full.html", layoutData)
}

func oldStockView(c *gin.Context) {

	var dispData []models.StockDataSaveFormat
	dispData = controllers.ShowOldStock("stock124")

	layoutData := struct {
		ThreadID int
		Posts    []models.StockDataSaveFormat
	}{
		ThreadID: 1,
		Posts:    dispData,
	}

	c.HTML(http.StatusOK, "sold.html", layoutData)
}

func insertStock(c *gin.Context) {
	c.HTML(http.StatusOK, "stock.html", nil)
}

func submit(c *gin.Context) {
	controllers.InsertStock("stock123", c.PostForm("stock"), c.PostForm("price"), c.PostForm("qty"))
}

func sell(c *gin.Context) {
	//log.Println(c.PostForm("symbol"))
	//log.Println(c.PostForm("salesprice"))
	controllers.SellStock("stock123", "stock124", c.PostForm("symbol"), c.PostForm("salesprice"))
}

// index sets up a session for the current user and serves the index page
func loginHandler(c *gin.Context) {
	// This check prevents the "/" handler from handling all requests by default
	/**
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}*/

	// Create a state token to prevent request forgery and store it in the session
	// for later validation
	session, err := store.Get(r, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		// Ignore the initial session fetch error, as Get() always returns a
		// session, even if empty.
		//return &appError{err, "Error fetching session", 500}
	}
	state := models.RandomString(64)
	session.Values["state"] = state
	session.Save()

	stateURL := url.QueryEscape(session.Values["state"].(string))

	// Fill in the missing fields in index.html
	/**
	var data = struct {
		ApplicationName, ClientID, State string
	}{applicationName, clientID, stateURL}*/

	// Render and serve the HTML
	err = c.HTML(http.StatusOK, "stock.html", nil)
	if err != nil {
		log.Println("error rendering template:", err)
		return &appError{err, "Error rendering template", 500}
	}
	return nil
}

// connect exchanges the one-time authorization code for a token and stores the
// token in the session
func connect(c *gin.Context) {
	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user
	session, err := store.Get(r, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		return &appError{err, "Error fetching session", 500}
	}
	if c.Query("state") != session.Values["state"].(string) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}
	// Normally, the state is a one-time token; however, in this example, we want
	// the user to be able to connect and disconnect without reloading the page.
	// Thus, for demonstration, we don't implement this best practice.
	// session.Values["state"] = nil

	// Setup for fetching the code from the request payload
	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	code := string(x)

	accessToken, idToken, err := models.Exchange(c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	gplusID, err := decodeIdToken(idToken)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check if the user is already connected
	storedToken := session.Values["accessToken"]
	storedGPlusID := session.Values["gplusID"]
	if storedToken != nil && storedGPlusID == gplusID {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Store the access token in the session for later use
	session.Values["accessToken"] = accessToken
	session.Values["gplusID"] = gplusID
	session.Save(r, w)
	return nil
}
