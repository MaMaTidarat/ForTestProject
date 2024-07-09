package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	app.Post("/import-file", importFileHandler)

	fmt.Println("Server started at :3000")
	log.Fatal(app.Listen(":3000"))
}

func importFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer excelFile.Close()

	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	finalResult := make(map[string][]map[string]int)

	for _, row := range rows[1:] { // Skipping header row
		dataCell := row[0]
		factorCell := row[1]

		if strings.Contains(dataCell, "-") {
			parts := strings.Split(dataCell, "-")
			from, _ := strconv.Atoi(parts[0])
			to, _ := strconv.Atoi(parts[1])

			fieldName := strings.ToLower(factorCell) // Convert Factor to lowercase (e.g., "CC" to "cc" and "AGE" to "age")
			entry := map[string]int{
				"from": from,
				"to":   to,
			}

			finalResult[fieldName] = append(finalResult[fieldName], entry)
		}
	}

	// Print the JSON data before saving
	for _, document := range finalResult {
		jsonData, err := json.MarshalIndent(document, "", "  ")
		if err != nil {
			return err
		}
		log.Println(string(jsonData))
	}

	// MongoDB connection
	mongoURI := "mongodb+srv://MaMa:AbpwIdEbqNsDcuks@cluster0.oakoge4.mongodb.net/"

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		}
	}()

	collection := client.Database("GI").Collection("product3")

	// Insert the final result into MongoDB
	_, err = collection.InsertOne(context.TODO(), finalResult)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("Data inserted successfully")
}
