package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var config Config

// Config is for the configuration of the database and other settings
type Config struct {
	DBURL                   string `json:"db_url"`
	DBName                  string `json:"db_name"`
	TrackDBCollectionName   string `json:"track_db_collection_name"`
	WebhookDBCollectionName string `json:"webhook_db_collection_name"`
	MaxTracksPerPage        int    `json:"max_tracks_per_page"`
}

func configure() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	startDatabases()
}

func startDatabases() {
	// Start a connection to our track global database
	trackGlobalDB = TrackMongoDB{
		config.DBURL,
		config.DBName,
		config.TrackDBCollectionName,
	}
	trackGlobalDB = trackGlobalDB.Init()

	// Start a connection to our webhook global database
	webhookGlobalDB = WebhookMongoDB{
		config.DBURL,
		config.DBName,
		config.WebhookDBCollectionName,
	}
	webhookGlobalDB = webhookGlobalDB.Init()
}
