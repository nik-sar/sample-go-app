package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/rand"
	"net/http"
	"sample/client"
	"sample/config"
	"sample/utils"
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

var storage *client.StorageType

const ShortUrlLength = 5

var cfg *config.AppConfig

func main() {
	cfg = config.GetAppConfig()
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	storage = client.CreateStorage(cfg, context.TODO())

	router := gin.Default()
	router.POST("/createShortUrl", postShortUrl)
	router.GET("/:urlParam", getLongUrlByShort)
	err := router.Run(":8181")
	if err != nil {
		return
	}
}

func getLongUrlByShort(c *gin.Context) {
	shortUrl := c.Param("urlParam")
	var link LinkModel
	err := client.FindOne(storage, bson.M{"shorturl": shortUrl}, &link)
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

	shortUrl := utils.GenerateRandomString(ShortUrlLength, random)
	//TODO: check unique shortUrl
	data := LinkModel{Url: request.Url, ShortUrl: shortUrl}
	err := client.Add(storage, data)
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	resp := LinkResponse{ShortUrl: cfg.Hostname + shortUrl}
	c.IndentedJSON(http.StatusCreated, resp)
}
