package models

//"github.com/fabwi987/go-getting-started/models"

//import "github.com/fabwi987/dbtest/models"

//StockDataSaveFormat is the format for saving the stock data
type StockDataSaveFormat struct {
	ID                 int     `json:"ID" bson:"ID"`
	Name               string  `json:"Name" bson:"Name"`
	Symbol             string  `json:"Symbol" bson:"Symbol"`
	LastTradePriceOnly float64 `json:"LastTradePriceOnly" bson:"LastTradePriceOnly"`
	Change             float64 `json:"Change" bson:"Change"`
	BuyPrice           float64 `json:"BuyPrice" bson:"BuyPrice"`
	NumberOfShares     float64 `json:"NumberOfShares" bson:"NumberOfShares"`
	Created            string  `json:"created" bson:"Created"`
	Progress           float64 `json:"Progress" bson:"Progress"`
	Color              string  `json:"Color" bson:"Color"`
	Active             string  `json:"Active" bson:"Active"`
	SalesPrice         float64 `json:"SalesPrice" bson:"SalesPrice"`
	Updated            string  `json:"Updated" bson:"Updated"`
	User               string  `json:"Updated" bson:"Updated"`
}

//StockDataDisplay format associates the stock with a user
type StockDataDisplayFormat struct {
	Stock     StockDataSaveFormat `json:"Stock" bson:"Stock"`
	StockUser User                `json:"StockUser" bson:"StockUser"`
}

//Datatype for diplaying stocks and users
type DispStocks []StockDataDisplayFormat
