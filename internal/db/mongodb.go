package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var MongoClient *mongo.Client

func InitMongoDB() {

	//creat a context with a timeout (cancel the context if the operation takes too long)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//establish a connection to the MongoDB server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://finbud123:finbud123@cluster0.8mbj0ln.mongodb.net/development?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	//Test ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	
	MongoClient = client

	log.Println("MongoDB connected successfully")
}
