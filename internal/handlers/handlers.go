package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"stock-update-lambda/internal/models"

	"stock-update-lambda/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockDataEntry struct {
	Symbol       string
	CurrentPrice float64 `json:"c"`
}

type StockData struct {
	symbol StockDataEntry
}

func FetchStockData(symbol string) (StockDataEntry, error) {
	var stockApiKey = "curorapr01qt2ncgdjvgcurorapr01qt2ncgdk00"
	resp, err := http.Get("https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=" + stockApiKey)

	if err != nil {
		log.Println("Error fetching stock data: ", err)
		return StockDataEntry{}, err
	}

	defer resp.Body.Close()

	var parsed struct {
		C float64 `json:"c"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Println("Error decoding stock data:", err)
		return StockDataEntry{}, err
	}

	stockData := StockDataEntry{
		CurrentPrice: parsed.C,
		Symbol:       symbol,
	}

	fmt.Println("Stock Data: ", stockData)
	log.Println("Fetched stock data successfully", stockData)
	return stockData, nil
}

func GetAllUserStocks() ([]string, []string, error) {
	collection := db.MongoClient.Database("development").Collection("userholdings")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching user stocks: ", err)
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var userIds []string
	var userHoldings []models.UserHolding
	for cursor.Next(ctx) {
		var holding models.UserHolding
		err := cursor.Decode(&holding)
		if err != nil {
			log.Println("Error decoding user stocks: ", err)
			return nil, nil, err
		}
		userHoldings = append(userHoldings, holding)
	}

	uniqueStocks := make(map[string]struct{})

	for _, holding := range userHoldings {
		userIds = append(userIds, holding.UserID.Hex())
		for _, stock := range holding.Holdings {
			uniqueStocks[stock.Symbol] = struct{}{}
		}
	}

	var result []string
	for stock := range uniqueStocks {
		result = append(result, stock)
	}

	// fmt.Println("Unique Stocks: ", result)
	log.Println("Fetched unique stocks successfully")
	return result, userIds, nil
}

func UpdatePortfolio(userID string, stockPrices models.GlobalStock) error {
	fmt.Println("Updating portfolio for user ID:", userID)
	holding_collection := db.MongoClient.Database("development").Collection("userholdings")
	user_collection := db.MongoClient.Database("development").Collection("users")
	portfolio_collection := db.MongoClient.Database("development").Collection("portfolios")

	// Fetch stock data for the user
	ctx, cancel := context.WithTimeout(context.Background(), 15*60*time.Second)
	defer cancel()

	var err error

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Invalid ObjectID format:", err)
		return err
	}

	var userHolding models.UserHolding
	err = holding_collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&userHolding)

	if err != nil {
		log.Println("Error fetching user holdings: ", err)
		return err
	}

	if len(userHolding.Holdings) == 0 {
		log.Println("No holdings found for user ID:", userID)
		return nil
	}

	var user models.User
	err = user_collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		log.Println("Error fetching user data: ", err)
		return err
	}

	var portfolio models.Portfolio
	err = portfolio_collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&portfolio)
	if err != nil {
		log.Println("Error fetching user portfolio: ", err)
		return err
	}

	var totalValue float64 = user.BankingAccountData.CashValue
	for _, stock := range userHolding.Holdings {
		updatedPrice, ok := stockPrices.Prices[stock.Symbol]
		if !ok {
			log.Printf("No price found for symbol %s", stock.Symbol)
			continue
		}
		totalValue += stock.Quantity * updatedPrice
	}

	user.BankingAccountData.StockValue = totalValue - user.BankingAccountData.CashValue
	user.BankingAccountData.AccountBalance = user.BankingAccountData.StockValue + user.BankingAccountData.CashValue

	_, err = user_collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": user})
	if err != nil {
		log.Println("Error updating user data: ", err)
	}

	today := time.Now()
	portfolioEntry := models.PortfolioEntry{
		Date:       today,
		TotalValue: totalValue,
	}

	var lastPortfolioEntry models.PortfolioEntry
	lastPortfolioEntry = portfolio.Portfolio[len(portfolio.Portfolio)-1]

	lastEntryDate := lastPortfolioEntry.Date
	if lastEntryDate.Year() == today.Year() &&
		lastEntryDate.Month() == today.Month() &&
		lastEntryDate.Day() == today.Day() {
		log.Println("Portfolio entry already exists for today")
		return nil
	}

	portfolio.Portfolio = append(portfolio.Portfolio, portfolioEntry)

	_, err = portfolio_collection.UpdateOne(ctx, bson.M{"userId": objectID}, bson.M{"$set": portfolio})

	if err != nil {
		log.Println("Error updating portfolio data: ", err)
		return err
	}
	log.Println("Portfolio updated successfully for user ID:", userID)
	return nil
}
