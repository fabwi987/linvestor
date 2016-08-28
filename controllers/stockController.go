package controllers

import (
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/fabwi987/linvestor/models"
)

//ShowStocks collects all stock from the database and gets their latest information from the finance api
func ShowStocks(dbtable string) ([]models.StockDataSaveFormat, string) {

	dbData, err := DBQuerySQL(dbtable)
	Perror(err)

	var allData = make(Stocks, len(dbData))
	var currString string

	var symbols []string
	for i := 0; i < len(allData); i++ {
		symbols[i] = allData[i].Symbol
	}

	for i := 0; i < len(allData); i++ {
		log.Println("Hämtar från api")
		//var data models.StockDataSaveFormat
		var moddatan models.StockDataSaveFormat
		//datan, err := yagoo.Get(dbData[i].Symbol)
		//moddatan, err := ModifyStock(datan, dbData[i].BuyPrice, dbData[i].LastTradePriceOnly)
		Perror(err)
		allData[i] = UpdateStock(moddatan, "current")
	}
	sort.Sort(allData)
	currString = ""
	return allData, currString
}

//UpdateStock takes a stock and updates with the latest data from the api and other minor operations
func UpdateStock(stock models.StockDataSaveFormat, kind string) models.StockDataSaveFormat {
	log.Println("Uppdaterar data")
	//buyPrice, err := strconv.ParseFloat(stock.BuyPrice, 64)
	//Perror(err)

	var dev float64
	if kind == "old" {
		//salesPrice, err := strconv.ParseFloat(stock.SalesPrice, 64)
		dev = RoundUp(stock.SalesPrice/stock.BuyPrice, 4)
		//Perror(err
	} else {
		//lastPrice, err := strconv.ParseFloat(stock.LastTradePriceOnly, 64)
		dev = RoundUp(stock.LastTradePriceOnly/stock.BuyPrice, 4)
		//Perror(err)
	}

	stock.Progress = (dev * 100) - 100
	//devString := strconv.FormatFloat(dev, 'f', 2, 64)
	//stock.Progress = devString + " %"

	if stock.Progress > 0 {
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
	//var stock models.StockDataSaveFormat
	//stock, err := yagoo.Get(_symbol)

	//price, err := strconv.ParseFloat(_price, 64)
	//number, err := strconv.ParseFloat(_number, 64)
	//Perror(err)

	//stockSave, err = ModifyStock(stock, price, number)
	//Perror(err)
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
	dbStock.SalesPrice, err = strconv.ParseFloat(price, 64)

	//Insert post to old stock db
	res, err := DbInsertSQL(dbStock, dbtableold)
	Perror(err)
	log.Println(res)

	//Delete post from active db
	err = DBDeletePost(dbtable, dbStock.Symbol)
}

//ModifyStock modifies the stock before insert to DB
/**
func ModifyStock(_stock yagoo.StockData, _buyPrice float64, _numberOfShares float64) (models.StockDataSaveFormat, error) {
	var toDB models.StockDataSaveFormat
	var err error
	toDB.Name = _stock.Query.Results.Quote.Name
	toDB.Symbol = _stock.Query.Results.Quote.Symbol
	toDB.LastTradePriceOnly, err = strconv.ParseFloat(_stock.Query.Results.Quote.LastTradePriceOnly, 64)
	toDB.Change, err = strconv.ParseFloat(_stock.Query.Results.Quote.Change, 64)
	toDB.Created = _stock.Query.Created
	toDB.BuyPrice = _buyPrice
	toDB.NumberOfShares = _numberOfShares
	toDB.Updated = _stock.Query.Created
	Perror(err)
	return toDB, nil
}*/

//Stocks is a slice of StockDataSaveFormat
type Stocks []models.StockDataSaveFormat

func (slice Stocks) Len() int { return len(slice) }
func (slice Stocks) Less(i, j int) bool {
	return (slice[i].LastTradePriceOnly / slice[i].BuyPrice) > (slice[j].LastTradePriceOnly / slice[j].BuyPrice)
}
func (slice Stocks) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}
