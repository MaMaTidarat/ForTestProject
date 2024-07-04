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

// MongoDB client
var client *mongo.Client

func main() {
	// Connect to MongoDB

	mongoURI := "mongodb+srv://MaMa:AbpwIdEbqNsDcuks@cluster0.oakoge4.mongodb.net/"
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Disconnect from MongoDB
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection closed.")

	// Create a new Fiber app
	app := fiber.New()

	// Route to upload Excel file
	app.Post("/import-file", importFileHandler)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func importFileHandler(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No file uploaded")
	}
	file := files[0]

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the Excel file
	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		return err
	}

	// Get rows from the specified sheet
	rows, err := excelFile.GetRows("Sheet1") // replace with your sheet name
	if err != nil {
		return err
	}

	var documents []map[string]interface{}

	// Parse the data from the columns
	for _, row := range rows {

		for _, cell := range row {
			log.Println("===>", string(cell))

			if strings.Contains(cell, "-") {
				parts := strings.Split(cell, "-")
				form, _ := strconv.Atoi(parts[0])
				to, _ := strconv.Atoi(parts[1])
				fieldName := "data" // replace this with the dynamic field name from your Excel column, e.g., "CC"
				doc := map[string]interface{}{
					fieldName: map[string]interface{}{
						"form": form,
						"to":   to,
					},
				}
				documents = append(documents, doc)
			}
		}
	}

	// Print the JSON data before saving
	for _, document := range documents {
		jsonData, err := json.MarshalIndent(document, "", "  ")
		if err != nil {
			return err
		}
		log.Println(string(jsonData))
	}

	// Insert the parsed data into MongoDB
	err = insertIntoMongoDB(documents)
	if err != nil {
		return err
	}

	return c.SendString("Data inserted successfully!")
}

func insertIntoMongoDB(documents []map[string]interface{}) error {
	collection := client.Database("GI").Collection("product3")

	for _, document := range documents {
		_, err := collection.InsertOne(context.Background(), document)
		if err != nil {
			return err
		}
	}

	fmt.Println("Data inserted successfully!")
	return nil
}
