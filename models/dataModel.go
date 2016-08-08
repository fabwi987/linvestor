package models

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
