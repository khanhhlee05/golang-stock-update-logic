package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"stock-update-lambda/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetAllUserStocks(collection *mongo.Collection) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching user stocks: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var userHoldings []models.UserHolding
	for cursor.Next(ctx) {
		var holding models.UserHolding
		err := cursor.Decode(&holding)
		if err != nil {
			log.Println("Error decoding user stocks: ", err)
			return nil, err
		}
		userHoldings = append(userHoldings, holding)
	}

	uniqueStocks := make(map[string]struct{})

	for _, holding := range userHoldings {
		for _, stock := range holding.Holdings {
			uniqueStocks[stock.Symbol] = struct{}{}
		}
	}

	var result []string
	for stock := range uniqueStocks {
		result = append(result, stock)
	}

	fmt.Println("Unique Stocks: ", result)
	return result, nil
}
