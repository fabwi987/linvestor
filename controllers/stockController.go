package controllers

import (
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/fabwi987/linvestor/models"
	"github.com/fabwi987/yaho"
)

//ShowStocks collects all stock from the database and gets their latest information from the finance api
func ShowStocks(dbtable string) ([]models.StockDataDisplayFormat, string) {

	dbData, err := DBQuerySQL(dbtable)
	Perror(err)

	var allData = make(Stocks, len(dbData))
	var currString string
	var symbols string

	//log.Println(len(allData))
	for i := 0; i < len(allData); i++ {
		symbols = symbols + "," + dbData[i].Symbol
		//log.Println("Hämtar från api")

	}
	datan, err := yaho.PolyGet(symbols)

	for i := 0; i < len(allData); i++ {

		moddatan, err := ModifyStock(datan.Query.Results.Quote[i], dbData[i].BuyPrice, dbData[i].LastTradePriceOnly, dbData[i].ID)
		Perror(err)
		allData[i] = UpdateStock(moddatan, "current")

	}
	sort.Sort(allData)
	usersArray := DBQueryUserRelation("usrtbl", allData)
	dispArray := make(models.DispStocks, len(allData))

	for i := 0; i < len(allData); i++ {
		dispArray[i].Stock = allData[i]
		dispArray[i].StockUser = usersArray[i]
	}
	currString = ""
	return dispArray, currString
}

//UpdateStock takes a stock and updates with the latest data from the api and other minor operations
func UpdateStock(stock models.StockDataSaveFormat, kind string) models.StockDataSaveFormat {
	log.Println("Uppdaterar data")
	//buyPrice, err := strconv.ParseFloat(stock.BuyPrice, 64)
	//Perror(err)

	var dev float64

	if kind == "old" {
		//salesPrice, err := strconv.ParseFloat(stock.SalesPrice, 64)
		//dev = RoundUp(stock.SalesPrice/stock.BuyPrice, 4)
		dev = stock.SalesPrice / stock.BuyPrice
		//dev = toFixed(stock.SalesPrice/stock.BuyPrice, 3)
		//Perror(err
	} else {
		//lastPrice, err := strconv.ParseFloat(stock.LastTradePriceOnly, 64)
		//dev = RoundUp(stock.LastTradePriceOnly/stock.BuyPrice, 4)
		dev = stock.LastTradePriceOnly / stock.BuyPrice
		//dev = toFixed(stock.LastTradePriceOnly/stock.BuyPrice, 3)
		//Perror(err)
	}

	dev = (dev * 100) - 100
	stock.ProgFloat = dev

	//log.Println(dev)

	if dev > 0 {
		stock.Color = "green"
	} else {
		stock.Color = "red"
	}

	procString := strconv.FormatFloat(dev, 'f', 2, 64)
	subString := procString
	//subString := procString[:4]
	stock.Progress = subString

	//devString := strconv.FormatFloat(dev, 'f', 2, 64)
	//stock.Progress = devString + " %"

	return stock

}

//ShowOldStock fetch stocks that have been sold
func ShowOldStock(dbtable string) []models.StockDataDisplayFormat {

	var dbData []models.StockDataSaveFormat
	dbData, err := DBQuerySQL(dbtable)
	Perror(err)

	usersArray := DBQueryUserRelation("usrtbl", dbData)
	dispArray := make(models.DispStocks, len(dbData))

	//var allData = make([]models.StockDataSaveFormat, len(dbData))

	for i := 0; i < len(dispArray); i++ {
		dispArray[i].Stock = UpdateStock(dbData[i], "old")
		dispArray[i].StockUser = usersArray[i]
	}

	return dispArray

}

//InsertStock Inserts the selected stock to the database
func InsertStock(dbtable string, _symbol string, _price string, _number string, _name string, prospect bool) {

	var stockSave models.StockDataSaveFormat
	stock, err := yaho.Get(_symbol)

	if prospect {
		_price = stock.Query.Results.Quote.LastTradePriceOnly
	}

	price, err := strconv.ParseFloat(_price, 64)
	number, err := strconv.ParseFloat(_number, 64)
	Perror(err)

	stockSave, err = ModifyStock(stock.Query.Results.Quote, price, number, 0)
	Perror(err)
	res, err := DbInsertSQL(stockSave, dbtable)
	Perror(err)

	//log.Println(res)
	//log.Println(_name)

	userid := DBQueryUserOnName(_name, "usrtbl")
	//log.Println(userid)
	DBCreateRelation(userid, res, "userrelation")

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
func ModifyStock(_stock yaho.Quote, _buyPrice float64, _numberOfShares float64, ID int) (models.StockDataSaveFormat, error) {
	var toDB models.StockDataSaveFormat
	var err error
	toDB.ID = ID
	toDB.Name = _stock.Name
	toDB.Symbol = _stock.Symbol
	toDB.LastTradePriceOnly, err = strconv.ParseFloat(_stock.LastTradePriceOnly, 64)
	toDB.Change, err = strconv.ParseFloat(_stock.Change, 64)
	toDB.Created = "0"
	toDB.BuyPrice = _buyPrice
	toDB.NumberOfShares = _numberOfShares
	toDB.Updated = "0"
	toDB.User = "0"
	Perror(err)
	return toDB, nil
}

//Stocks is a slice of StockDataSaveFormat
type Stocks []models.StockDataSaveFormat

//Stocks is a slice of StockDataSaveFormat
type StocksDisp []models.StockDataDisplayFormat

func (slice Stocks) Len() int { return len(slice) }
func (slice Stocks) Less(i, j int) bool {
	return slice[i].ProgFloat > slice[j].ProgFloat
}
func (slice Stocks) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

//RoundUp removes a number of decimals
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

//round rounds a number
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

//toFixed sets the number of decimals
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
