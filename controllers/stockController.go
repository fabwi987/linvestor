package controllers

import (
	"log"
	"strconv"

	"github.com/fabwi987/YaGoo"
	"github.com/fabwi987/linvestor/models"
)

//ShowStocks collects all stock from the database and gets their latest information from the finance api
func ShowStocks(dbtable string) ([]models.StockDataSaveFormat, string) {

	dbData, err := DBQuerySQL(dbtable)
	Perror(err)

	var allData = make([]models.StockDataSaveFormat, len(dbData))
	var currString string

	for i := 0; i < len(allData); i++ {

		datan, err := yagoo.Get(dbData[i].Symbol)
		moddatan, err := ModifyStock(datan, dbData[i].BuyPrice, dbData[i].LastTradePriceOnly)
		Perror(err)
		allData[i] = UpdateStock(moddatan, "current")
	}

	currString = ""
	return allData, currString
}

//UpdateStock takes a stock and updates with the latest data from the api and other minor operations
func UpdateStock(stock models.StockDataSaveFormat, kind string) models.StockDataSaveFormat {

	buyPrice, err := strconv.ParseFloat(stock.BuyPrice, 64)
	Perror(err)

	var dev float64
	if kind == "old" {
		salesPrice, err := strconv.ParseFloat(stock.SalesPrice, 64)
		dev = salesPrice / buyPrice
		Perror(err)
	} else {
		lastPrice, err := strconv.ParseFloat(stock.LastTradePriceOnly, 64)
		dev = lastPrice / buyPrice
		Perror(err)
	}

	dev = (dev * 100) - 100
	devString := strconv.FormatFloat(dev, 'f', 2, 64)
	stock.Progress = devString + " %"

	if dev > 0 {
		stock.Color = "green"
	} else {
		stock.Color = "red"
	}

	return stock

}

//ShowOldStock fetch stocks that have been sold
func ShowOldStock(dbtable string) []models.StockDataSaveFormat {

	var dbData []models.StockDataSaveFormat
	dbData, err := DBQuerySQL(dbtable)
	Perror(err)

	var allData = make([]models.StockDataSaveFormat, len(dbData))

	for i := 0; i < len(allData); i++ {
		allData[i] = UpdateStock(dbData[i], "old")
	}

	return allData

}

//InsertStock Inserts the selected stock to the database
func InsertStock(dbtable string, _symbol string, _price string, _number string) {

	var stockSave models.StockDataSaveFormat
	stock, err := yagoo.Get(_symbol)

	stockSave, err = ModifyStock(stock, _price, _number)
	Perror(err)
	res, err := DbInsertSQL(stockSave, dbtable)
	Perror(err)
	log.Println(res)
}

//SellStock takes the sales price as input and moves the stock to the old stock table
func SellStock(dbtable string, dbtableold string, symbol string, price string) {

	log.Println("Symbol = " + symbol)
	log.Println("Sales price = " + price)
	//Create temp variable to hold stock

	var dbStock models.StockDataSaveFormat

	//Find stock data from db
	dbStock, err := DBQuerySQLSingle(symbol, dbtable)
	log.Println(dbStock.Symbol)

	//Add the sales pric from user input
	dbStock.SalesPrice = price

	//Insert post to old stock db
	res, err := DbInsertSQL(dbStock, dbtableold)
	Perror(err)
	log.Println(res)

	//Delete post from active db
	err = DBDeletePost(dbtable, dbStock.Symbol)
}

//ModifyStock modifies the stock before insert to DB
func ModifyStock(_stock yagoo.StockData, _buyPrice string, _numberOfShares string) (models.StockDataSaveFormat, error) {
	var toDB models.StockDataSaveFormat
	toDB.Name = _stock.Query.Results.Quote.Name
	toDB.Symbol = _stock.Query.Results.Quote.Symbol
	toDB.LastTradePriceOnly = _stock.Query.Results.Quote.LastTradePriceOnly
	toDB.Change = _stock.Query.Results.Quote.Change
	toDB.Created = _stock.Query.Created
	toDB.BuyPrice = _buyPrice
	toDB.NumberOfShares = _numberOfShares
	toDB.Updated = _stock.Query.Created
	return toDB, nil
}
