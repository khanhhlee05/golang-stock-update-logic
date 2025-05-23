package main

import (
	"fmt"
	"stock-update-lambda/internal/db"
	"stock-update-lambda/internal/handlers"
	"stock-update-lambda/internal/models"
	"sync"
	"time"
)

func measureExecutionTime(name string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		fmt.Printf("%s took %v\n", name, duration)
	}
}

func main() {
	defer measureExecutionTime("Main Execution")()
	db.InitMongoDB()

	fmt.Println("MongoDB Client Initialized: ", db.MongoClient != nil)

	// Initialize as a map instead of GlobalStock
	stockPrices := make(map[string]float64)

	uniqueStock, userIds, err := handlers.GetAllUserStocks()
	if err != nil {
		fmt.Println("Error fetching unique stocks: ", err)
		return
	}

	for _, stock := range uniqueStock {
		current, err := handlers.FetchStockData(stock)
		if err != nil {
			fmt.Println("Error fetching stock data: ", err)
			return
		}
		stockPrices[stock] = current.CurrentPrice
		time.Sleep(1 * time.Second) // Sleep for 1 second to avoid hitting the API rate limit
	}

	globalStocks := models.GlobalStock{
		Prices: stockPrices,
	}

	// Update the portfolio for each user
	// using goroutines for concurrent processing

	var wg sync.WaitGroup

	errorCh := make(chan error, len(userIds))

	for _, userId := range userIds {
		wg.Add(1)
		go func(uid string) {
			defer wg.Done()
			err := handlers.UpdatePortfolio(uid, globalStocks)
			if err != nil {
				errorCh <- err
			}
		}(userId)
	}

	wg.Wait()

	close(errorCh)

	for err := range errorCh {
		fmt.Println("Error updating portfolio:", err)
	}
}
