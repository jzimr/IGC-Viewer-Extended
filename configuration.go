package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var config Config

// Config is for the configuration of the database and other settings
type Config struct {
	DBURL            string `json:"db_url"`
	DBName           string `json:"db_name"`
	DBCollectionName string `json:"db_collection_name"`
	MaxTracksPerPage int    `json:"max_tracks_per_page"`
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
}
