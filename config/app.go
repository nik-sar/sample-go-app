package config

import (
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	Hostname            string
	MongoUri            string
	MongoDbName         string
	MongoCollectionName string
}

func GetAppConfig() *AppConfig {

	hostname := os.Getenv("HOSTNAME")
	if len(strings.TrimSpace(hostname)) == 0 {
		log.Fatal("HOSTNAME env is required")
	}
	mongoUri := os.Getenv("MONGODB_CONNECTION_URI")
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

	appConfig := AppConfig{
		Hostname:            hostname,
		MongoUri:            mongoUri,
		MongoDbName:         mongoDbName,
		MongoCollectionName: mongoCollectionName,
	}

	return &appConfig
}
