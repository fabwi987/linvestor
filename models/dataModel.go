package models

import (
	"github.com/fabwi987/YaGoo"
)

//Modified struct for saving the data
type StockDataSaveFormat struct {
	Name               string `json:"Name" bson:"Name"`
	Symbol             string `json:"Symbol" bson:"Symbol"`
	LastTradePriceOnly string `json:"LastTradePriceOnly" bson:"LastTradePriceOnly"`
	Change             string `json:"Change" bson:"Change"`
	BuyPrice           string `json:"BuyPrice" bson:"BuyPrice"`
	NumberOfShares     string `json:"NumberOfShares" bson:"NumberOfShares"`
	Created            string `json:"created" bson:"Created"`
	Progress           string `json:"Progress" bson:"Progress"`
	Color              string `json:"Color" bson:"Color"`
	Active             string `json:"Active" bson:"Active"`
	SalesPrice         string `json:"SalesPrice" bson:"SalesPrice"`
}

//Modifies the stock before insert to DB
func ModifyStock(_stock yagoo.StockData, _buyPrice string, _numberOfShares string) (StockDataSaveFormat, error) {
	var toDB StockDataSaveFormat
	toDB.Name = _stock.Query.Results.Quote.Name
	toDB.Symbol = _stock.Query.Results.Quote.Symbol
	toDB.LastTradePriceOnly = _stock.Query.Results.Quote.LastTradePriceOnly
	toDB.Change = _stock.Query.Results.Quote.Change
	toDB.Created = _stock.Query.Created
	toDB.BuyPrice = _buyPrice
	toDB.NumberOfShares = _numberOfShares
	return toDB, nil
}
