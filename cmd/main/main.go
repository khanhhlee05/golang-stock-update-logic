package main

import (
	"fmt"
	"stock-update-lambda/internal/db"

)

func main() {
	db.InitMongoDB()

	fmt.Println("MongoDB Client Initialized: ", db.MongoClient != nil)

	
}