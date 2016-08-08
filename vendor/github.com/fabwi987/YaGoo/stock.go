package yagoo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

//Get calls the API and receives a json response
func Get(symbol string) (StockData, error) {

	//Build the http request
	query := url.QueryEscape("select * from yahoo.finance.quote where symbol in ('" + symbol + "')")
	destination := "http://query.yahooapis.com/v1/public/yql?q="
	specification := "&format=json&diagnostics=true&env=http://datatables.org/alltables.env"

	resp, err := http.Get(destination + query + specification)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	var stock StockData

	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		return stock, err
	}

	return stock, nil

}

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
