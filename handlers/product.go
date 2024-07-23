package main

import (
	"database"

	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProducts(c *fiber.Ctx) error {
	param := c.Query("param")
	status := c.Query("status")

	filter := bson.M{}
	if param != "" {
		filter["$or"] = []bson.M{
			{"subProductGroup.key": bson.M{"$regex": param, "$options": "i"}},
			{"key": bson.M{"$regex": param, "$options": "i"}},
			{"productList.productName": bson.M{"$regex": param, "$options": "i"}},
			{"productList.insurer.insurerCode": bson.M{"$regex": param, "$options": "i"}},
			{"productList.brokers.key": bson.M{"$regex": param, "$options": "i"}},
		}
	}

	if status != "" {
		filter["productList.status"] = status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.ProductCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(products)
}
