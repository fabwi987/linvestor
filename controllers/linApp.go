package controllers

import (
	"log"
	"strconv"

	"github.com/fabwi987/linvestor/models"
)

//ShowStocks collects all stock from the database and gets their latest information from the finance api
func ShowStocks(dbtable string) ([]models.StockDataSaveFormat, string) {

	var dbData []models.StockDataSaveFormat
	dbData, err := models.DBQuerySQL(dbtable)
	models.Perror(err)

	var allData = make([]models.StockDataSaveFormat, len(dbData))
	//var startValue float64
	//var currentValue float64
	var currString string
	//For each stock, get latest value and update DB
	for i := 0; i < len(dbData); i++ {
		datan, err := models.GetShareData(dbData[i].Symbol)
		models.Perror(err)

		modDatan, err := models.ModifyStock(datan, dbData[i].BuyPrice, dbData[i].NumberOfShares)

		//Calcullate the change and set approperiate color
		lastPrice, err := strconv.ParseFloat(modDatan.LastTradePriceOnly, 64)
		buyPrice, err := strconv.ParseFloat(modDatan.BuyPrice, 64)
		//numberOfShares, err := strconv.ParseFloat(modDatan.NumberOfShares, 64)
		models.Perror(err)

		dev := lastPrice / buyPrice

		dev = (dev * 100) - 100
		devString := strconv.FormatFloat(dev, 'f', 2, 64)
		modDatan.Progress = devString + " %"

		if dev > 0 {
			modDatan.Color = "green"
		} else {
			modDatan.Color = "red"
		}

		//removed since we only display progress per stock
		//startValue = startValue + buyPrice*numberOfShares
		//currentValue = currentValue + lastPrice*numberOfShares

		models.Perror(err)
		allData[i] = modDatan
		//log.Println(i)
		//log.Println(allData[i].Symbol)
	}

	//Removed since we only show progress per stock
	//currentValue = currentValue / startValue
	//currentValue = (currentValue * 100) - 100
	//currString = strconv.FormatFloat(currentValue, 'f', 2, 64) + " %"
	currString = ""
	return allData, currString
}

//ShowOldStock fetch stocks that have been sold
func ShowOldStock(dbtable string) []models.StockDataSaveFormat {

	var dbData []models.StockDataSaveFormat
	dbData, err := models.DBQuerySQL(dbtable)
	models.Perror(err)

	var allData = make([]models.StockDataSaveFormat, len(dbData))

	//For each stock, get latest value and update DB
	for i := 0; i < len(dbData); i++ {
		//Calcullate the change and set approperiate color
		startValue, err := strconv.ParseFloat(dbData[i].BuyPrice, 64)
		saleValue, err := strconv.ParseFloat(dbData[i].SalesPrice, 64)
		models.Perror(err)

		dev := saleValue / startValue

		dev = (dev * 100) - 100
		devString := strconv.FormatFloat(dev, 'f', 2, 64)
		dbData[i].Progress = devString + " %"

		if dev > 0 {
			dbData[i].Color = "green"
		} else {
			dbData[i].Color = "red"
		}

		allData[i] = dbData[i]
	}

	return allData

}

//InsertStock Inserts the selected stock to the database
func InsertStock(dbtable string, _symbol string, _price string, _number string) {

	var stockSave models.StockDataSaveFormat
	stock, err := models.GetShareData(_symbol)

	stockSave, err = models.ModifyStock(stock, _price, _number)
	models.Perror(err)
	res, err := models.DbInsertSQL(stockSave, dbtable)
	models.Perror(err)
	log.Println(res)
}

//SellStock takes the sales price as input and moves the stock to the old stock table
func SellStock(dbtable string, dbtableold string, symbol string, price string) {

	log.Println("Symbol = " + symbol)
	log.Println("Sales price = " + price)
	//Create temp variable to hold stock

	var dbStock models.StockDataSaveFormat

	//Find stock data from db
	dbStock, err := models.DBQuerySQLSingle(symbol, dbtable)
	log.Println(dbStock.Symbol)

	//Add the sales pric from user input
	dbStock.SalesPrice = price

	//Insert post to old stock db
	res, err := models.DbInsertSQL(dbStock, dbtableold)
	models.Perror(err)
	log.Println(res)

	//Delete post from active db
	err = models.DBDeletePost(dbtable, dbStock.Symbol)
}
