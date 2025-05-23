package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StockData struct {
	symbol       string
	currentPrice float64 `json:"c"`
}

func FetchStockData(symbol string) (StockData, error) {
	var stockApiKey = "curorapr01qt2ncgdjvgcurorapr01qt2ncgdk00"
	resp, err := http.Get("https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=" + stockApiKey)

	if err != nil {
		log.Println("Error fetching stock data: ", err)
		return StockData{}, err
	}

	defer resp.Body.Close()

	var parsed struct {
		C float64 `json:"c"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Println("Error decoding stock data:", err)
		return StockData{}, err
	}

	stockData := StockData{
		currentPrice: parsed.C,
		symbol:       symbol,
	}

	fmt.Println("Stock Data: ", stockData)
	return stockData, nil
}
