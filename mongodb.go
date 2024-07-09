package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertIntoMongoDB(data map[string]interface{}) error {
	mongoURI := "mongodb+srv://MaMa:AbpwIdEbqNsDcuks@cluster0.oakoge4.mongodb.net/"

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		}
	}()

	collection := client.Database("GI").Collection("product3")

	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}
