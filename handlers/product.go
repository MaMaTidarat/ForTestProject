package handlers

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/MaMaTidarat/poc-app/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Product struct {
	ID           string       `json:"id"`
	ProductName  string       `json:"productName"`
	ProductGroup ProductGroup `json:"productGroup"`
	ProductType  ProductType  `json:"productType"`
	Insurer      Insurer      `json:"insurer"`
	Brokers      []Broker     `json:"brokers"`
	Status       string       `json:"status"`
}

type ProductGroup struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ProductType struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Insurer struct {
	ID          string `json:"_id"`
	InsurerCode string `json:"insurerCode"`
	InsurerName string `json:"insurerName"`
}

type Broker struct {
	Key         string `json:"key"`
	ChannelName string `json:"channelName"`
}

func SanitizeString(input string) string {
	re := regexp.MustCompile(`[.*+?^${}()|[\]\\]`)
	return re.ReplaceAllString(input, `\$0`)
}

func GetProducts(c *fiber.Ctx) error {
	param := c.Query("param")
	status := c.Query("status")

	log.Printf("Received param: %s", param)
	log.Printf("Received status: %s", status)

	filter := bson.M{}
	if param != "" {
		sanitizedParam := SanitizeString(param)
		log.Printf("Sanitized param: %s", sanitizedParam)
		filter["$or"] = []bson.M{
			{"productType.key": bson.M{"$regex": sanitizedParam, "$options": "i"}},
			{"key": bson.M{"$regex": sanitizedParam, "$options": "i"}},
			{"productList.productName": bson.M{"$regex": sanitizedParam, "$options": "i"}},
			{"productList.insurer.insurerCode": bson.M{"$regex": sanitizedParam, "$options": "i"}},
			{"productList.brokers.key": bson.M{"$regex": sanitizedParam, "$options": "i"}},
		}
		log.Printf("Filter created: %+v", filter)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.ProductCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding products: %v", err)
		return c.Status(500).SendString(err.Error())
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("Error decoding products: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	var products []Product
	for _, result := range results {
		productList, ok := result["productList"].(bson.A)
		if !ok {
			log.Printf("Error casting productList: %v", result["productList"])
			continue
		}

		for _, item := range productList {
			productMap, ok := item.(bson.M)
			if !ok {
				log.Printf("Error casting productMap: %v", item)
				continue
			}

			brokers := []Broker{}
			brokersData, ok := productMap["brokers"].(bson.A)
			if ok {
				for _, brokerItem := range brokersData {
					brokerMap, ok := brokerItem.(bson.M)
					if ok {
						brokers = append(brokers, Broker{
							Key:         getStringField(brokerMap, "key"),
							ChannelName: getStringField(brokerMap, "channelName"),
						})
					}
				}
			}

			product := Product{
				ID:          getStringField(productMap, "id"),
				ProductName: getStringField(productMap, "productName"),
				ProductGroup: ProductGroup{
					Name: getStringField(result, "name"),
					Key:  getStringField(result, "key"),
				},
				ProductType: ProductType{
					Name: getStringField(result["productType"].(bson.M), "name"),
					Key:  getStringField(result["productType"].(bson.M), "key"),
				},
				Insurer: Insurer{
					ID:          getStringField(productMap["insurer"].(bson.M), "_id"),
					InsurerCode: getStringField(productMap["insurer"].(bson.M), "insurerCode"),
					InsurerName: getStringField(productMap["insurer"].(bson.M), "insurerName"),
				},
				Brokers: brokers,
				Status:  getStringField(productMap, "productStatus"),
			}
			products = append(products, product)
		}
	}

	// If status parameter is provided, filter products by status
	if status != "" {
		sanitizedStatus := SanitizeString(strings.ToUpper(status))
		log.Printf("Sanitized status: %s", sanitizedStatus)
		var filteredProducts []Product
		for _, product := range products {
			if strings.EqualFold(product.Status, sanitizedStatus) {
				filteredProducts = append(filteredProducts, product)
			}
		}
		products = filteredProducts
	}

	return c.JSON(products)
}

// Helper function to safely get string field from a map
func getStringField(data interface{}, key string) string {
	if dataMap, ok := data.(bson.M); ok {
		if value, ok := dataMap[key].(string); ok {
			return value
		}
	}
	log.Printf("Field %s not found or not a string", key)
	return ""
}
