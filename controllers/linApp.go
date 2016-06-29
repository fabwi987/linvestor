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
	var startValue float64
	var currentValue float64
	var currString string
	//For each stock, get latest value and update DB
	for i := 0; i < len(dbData); i++ {
		datan, err := models.GetShareData(dbData[i].Symbol)
		models.Perror(err)

		modDatan, err := models.ModifyStock(datan, dbData[i].BuyPrice, dbData[i].NumberOfShares)

		//Calcullate the change and set approperiate color
		lastPrice, err := strconv.ParseFloat(modDatan.LastTradePriceOnly, 64)
		buyPrice, err := strconv.ParseFloat(modDatan.BuyPrice, 64)
		numberOfShares, err := strconv.ParseFloat(modDatan.NumberOfShares, 64)
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

		startValue = startValue + buyPrice*numberOfShares
		currentValue = currentValue + lastPrice*numberOfShares

		models.Perror(err)
		allData[i] = modDatan
		log.Println(i)
		log.Println(allData[i].Symbol)
	}

	currentValue = currentValue / startValue
	currentValue = (currentValue * 100) - 100
	currString = strconv.FormatFloat(currentValue, 'f', 2, 64) + " %"
	return allData, currString
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
