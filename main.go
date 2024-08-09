package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type LinkRequest struct {
	Url string `json:"url"`
}

type LinkResponse struct {
	ShortUrl string `json:"shortUrl"`
}

type LinkModel struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

var random *rand.Rand
var collection *mongo.Collection
var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const ShortUrlLength = 5

var hostname string
var mongoUri string

func main() {
	hostname = os.Getenv("HOSTNAME")
	if len(strings.TrimSpace(hostname)) == 0 {
		log.Fatal("HOSTNAME env is required")
	}
	mongoUri = os.Getenv("MONGODB_CONNECTION_URI")
	if len(strings.TrimSpace(mongoUri)) == 0 {
		log.Fatal("MONGODB_CONNECTION_URI env is required")
	}

	mongoDbName := os.Getenv("MONGODB_NAME")
	if len(strings.TrimSpace(mongoDbName)) == 0 {
		log.Fatal("MONGODB_NAME env is required")
	}

	mongoCollectionName := os.Getenv("MONGODB_COLLECTION")
	if len(strings.TrimSpace(mongoCollectionName)) == 0 {
		log.Fatal("MONGODB_COLLECTION env is required")
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	db := mongoConnect()
	collection = db.Database(mongoDbName).Collection(mongoCollectionName)

	router := gin.Default()
	router.POST("/createShortUrl", postShortUrl)
	router.GET("/:urlParam", getLongUrlByShort)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

func getLongUrlByShort(c *gin.Context) {
	shortUrl := c.Param("urlParam")
	var link LinkModel
	err := collection.FindOne(c, bson.M{"shorturl": shortUrl}).Decode(&link)
	if err != nil {
		log.Print(err)
		_ = c.AbortWithError(http.StatusNotFound, nil)
	}
	c.Redirect(http.StatusMovedPermanently, link.Url)
	//c.IndentedJSON(http.StatusOK, link)
}

func postShortUrl(c *gin.Context) {
	var request LinkRequest
	if err := c.BindJSON(&request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, nil)
		return
	}

	shortUrl := generateRandomString(ShortUrlLength)
	//TODO: check unique shortUrl
	data := LinkModel{Url: request.Url, ShortUrl: shortUrl}
	_, err := collection.InsertOne(c, data)
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	resp := LinkResponse{ShortUrl: hostname + shortUrl}
	c.IndentedJSON(http.StatusCreated, resp)
}

func generateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(b)
}

func mongoConnect() *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	return client
}
