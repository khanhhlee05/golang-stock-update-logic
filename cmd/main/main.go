package main

import (
	"fmt"
	"stock-update-lambda/internal/db"
	"stock-update-lambda/internal/handlers"
)

func main() {
	db.InitMongoDB()

	fmt.Println("MongoDB Client Initialized: ", db.MongoClient != nil)

	// handlers.FetchStockData("ggesttdadfasd")
	handlers.GetAllUserStocks(db.MongoClient.Database("development").Collection("userholdings"))
}
