package models

//Struct that is returned from the API
type StockData struct {
	Query struct {
		Count       int    `json:"count" bson:"count"`
		Created     string `json:"created" bson:"created"`
		Diagnostics struct {
			Build_version string `json:"build-version"`
			Cache         struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Method               string `json:"method"`
				Type                 string `json:"type"`
			} `json:"cache"`
			Javascript struct {
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Instructions_used    string `json:"instructions-used"`
				Table_name           string `json:"table-name"`
			} `json:"javascript"`
			PubliclyCallable string `json:"publiclyCallable"`
			Query            struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Params               string `json:"params"`
			} `json:"query"`
			Service_time string `json:"service-time"`
			URL          []struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
			} `json:"url"`
			User_time string `json:"user-time"`
		} `json:"diagnostics"`
		Lang    string `json:"lang"`
		Results struct {
			Quote struct {
				AverageDailyVolume   interface{} `json:"AverageDailyVolume" bson:"averagedailyvolume"`
				Change               string      `json:"Change" bson:"change"`
				DaysHigh             string      `json:"DaysHigh" bson:"dayshigh"`
				DaysLow              string      `json:"DaysLow" bson:"dayslow"`
				DaysRange            string      `json:"DaysRange" bson:"daysrange"`
				LastTradePriceOnly   string      `json:"LastTradePriceOnly" bson:"lasttradepriceonly"`
				MarketCapitalization string      `json:"MarketCapitalization" bson:"marketcapitalization"`
				Name                 string      `json:"Name" bson:"name" bson:"name"`
				StockExchange        string      `json:"StockExchange" bson:"stockexchange"`
				Symbol               string      `json:"Symbol" bson:"symbol"`
				Volume               string      `json:"Volume" bson:"volume"`
				YearHigh             string      `json:"YearHigh" bson:"yearhigh"`
				YearLow              string      `json:"YearLow" bson:"yearlow"`
				Symbol1              string      `json:"symbol" bson:"symbol1"`
			} `json:"quote" bson:"quote"`
		} `json:"results" bson:"results"`
	} `json:"query" bson:"query"`
}

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
func ModifyStock(_stock StockData, _buyPrice string, _numberOfShares string) (StockDataSaveFormat, error) {
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
