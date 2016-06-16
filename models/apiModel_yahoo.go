package models

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

//Calls the API and receives a json response
func GetShareData(_share string) (StockData, error) {
	query := url.QueryEscape("select * from yahoo.finance.quote where symbol in ('" + _share + "')")
	resp, err := http.Get("http://query.yahooapis.com/v1/public/yql?q=" + query + "&format=json&diagnostics=true&env=http://datatables.org/alltables.env")

	Perror(err)

	defer resp.Body.Close()

	var d StockData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Println(err)
	}

	log.Println("Hämtar data från API")
	return d, nil

}
