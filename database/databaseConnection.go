package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Printf("error loading .env file: %v", err)
	// }

	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func OpenCollection(databaseName, collectionName string) *mongo.Collection {
	client := DBInstance()

	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
