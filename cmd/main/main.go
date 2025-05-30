package main

import (
	"context"
	"fmt"
	"stock-update-lambda/internal/db"
	"stock-update-lambda/internal/handlers"
	"stock-update-lambda/internal/models"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func measureExecutionTime(name string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		fmt.Printf("%s took %v\n", name, duration)
	}
}

func handler(ctx context.Context) (string, error) {
	defer measureExecutionTime("Lambda Execution")()
	db.InitMongoDB()

	fmt.Println("MongoDB Client Initialized: ", db.MongoClient != nil)

	// Initialize as a map instead of GlobalStock
	stockPrices := make(map[string]float64)

	uniqueStock, userIds, err := handlers.GetAllUserStocks()
	if err != nil {
		fmt.Println("Error fetching unique stocks: ", err)
		return "Error fetching unique stocks", err
	}

	for _, stock := range uniqueStock {
		current, err := handlers.FetchStockData(stock)
		if err != nil {
			fmt.Println("Error fetching stock data: ", err)
			return "Error fetching stock data", err
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
		// It might be better to collect all errors and return them
		// For simplicity, returning the first error encountered for now
		return "Error updating portfolio", err
	}
	return "Successfully updated portfolios", nil
}

func main() {
	lambda.Start(handler)
}
