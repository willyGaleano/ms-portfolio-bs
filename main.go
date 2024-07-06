package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "ms-portfolio-bs/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set in .env file")
	}

	if err := connect_to_mongodb(uri); err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	println("Connected to MongoDB!")
}

// @title MS Portfolio BS API
// @description This is a simple API for managing portfolios
// @version 1.0
// @host localhost:3002
// @BasePath /ms-portfolio-bs/v1
func main() {
	r := gin.Default()

	// Use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.GET("ms-portfolio-bs/v1/portfolios/:id", getPortfolioByID)
	r.POST("ms-portfolio-bs/v1/portfolios/seed", seedData)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3002"
	}

	r.Run(":" + port)
}

func connect_to_mongodb(uri string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	mongoClient = client
	return err
}

// getPortfolioByID godoc
// @Summary Get portfolio by ID
// @Description get portfolio by ID
// @Tags portfolio
// @Accept  json
// @Produce  json
// @Param id path string true "Portfolio ID"
// @Success 200 {object} map[string]interface{}
// @Router /portfolios/{id} [get]
func getPortfolioByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var portfolio bson.M
	err = mongoClient.Database("portfolio_db").Collection("portfolio").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&portfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "OK",
		"data": portfolio,
	})
}

// seedData godoc
// @Summary Seed data into MongoDB
// @Description Seed data into MongoDB
// @Tags portfolio
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /portfolios/seed [post]
func seedData(c *gin.Context) {
	fileBytes, err := os.ReadFile("client_portfolio.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var documents []interface{}
	err = json.Unmarshal(fileBytes, &documents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, doc := range documents {
		docMap := doc.(map[string]interface{})
		if oid, ok := docMap["_id"].(map[string]interface{}); ok && oid["$oid"] != nil {
			oidStr := oid["$oid"].(string)
			objectID, err := primitive.ObjectIDFromHex(oidStr)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ObjectID"})
				return
			}
			docMap["_id"] = objectID
		}
		if createdDate, ok := docMap["createdDate"].(map[string]interface{}); ok && createdDate["$date"] != nil {
			docMap["createdDate"] = createdDate["$date"]
		}
	}

	collection := mongoClient.Database("portfolio_db").Collection("portfolio")
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data successfully seeded into MongoDB"})
}
