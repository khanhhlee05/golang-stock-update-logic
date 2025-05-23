package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Portfolio struct {
	UserID    primitive.ObjectID `bson:"userId"`
	Portfolio []PortfolioEntry   `bson:"portfolio"`
}

type PortfolioEntry struct {
	Date       time.Time `bson:"date"`
	TotalValue float64   `bson:"totalValue"`
}

type UserHolding struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"userId"`
	Holdings []UserHoldingEntry `bson:"stocks"`
}

type UserHoldingEntry struct {
	Symbol        string  `bson:"stockSymbol"`
	Quantity      float64 `bson:"quantity"`
	PurchasePrice float64 `bson:"purchasePrice"`
}

type User struct {
	BankingAccountData BankingAccountData `bson:"bankingAccountData"`
}

type BankingAccountData struct {
	AccountBalance float64 `bson:"accountBalance"`
	StockValue     float64 `bson:"stockValue"`
	CashValue      float64 `bson:"cash"`
}

type GlobalStock struct {
	Prices map[string]float64
}

//TODO: presave method: accountBalance = stockValue + cashValue
