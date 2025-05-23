package main

import (
	"fmt"
	"stock-update-lambda/internal/db"
	"stock-update-lambda/internal/handlers"
	"stock-update-lambda/internal/models"
)

func main() {
	db.InitMongoDB()

	fmt.Println("MongoDB Client Initialized: ", db.MongoClient != nil)

	var stockPrices models.GlobalStock
	// handlers.FetchStockData("ggesttdadfasd")
	uniqueStock, userIds, err := handlers.GetAllUserStocks()
	if err != nil {
		fmt.Println("Error fetching unique stocks: ", err)
		return
	}
	fmt.Println("Unique Stocks: ", uniqueStock)
	fmt.Println("User IDs: ", userIds)

	handlers.UpdatePortfolio(userIds[0], stockPrices)
}
