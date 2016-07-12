package main

import (
	"net/http"
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
	router.GET("/auth", models.AuthenticationHandler)
	router.GET("/connect", models.Connect)
	router.GET("/loggedin", models.Loggedin)
	router.POST("/disconnect", models.Disconnect)
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
