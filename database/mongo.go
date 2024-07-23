package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProductCollection *mongo.Collection

func ConnectDB() {
	mongoURI := "mongodb+srv://MaMa:AbpwIdEbqNsDcuks@cluster0.oakoge4.mongodb.net/"

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ProductCollection = client.Database("GI").Collection("productV4")
	log.Println("Connected to MongoDB!")
}
