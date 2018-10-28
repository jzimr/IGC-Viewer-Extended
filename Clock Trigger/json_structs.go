package main

// MongoDB stores the information of the DB connection
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

// PostSlackWebhook contains the message we want to send to the discord webhook
type PostSlackWebhook struct {
	Text string `json:"text"`
}

// TrackMetaData includes meta information about a particular track
type TrackMetaData struct {
	ID          string  `json:"id"`
	Hdate       string  `json:"H_date"`
	Pilot       string  `json:"pilot"`
	Glider      string  `json:"glider"`
	GliderID    string  `json:"glider_id"`
	TrackLength float64 `json:"track_length"`
	TrackSrcURL string  `json:"track_src_url"`
	Timestamp   int64   `json:"Timestamp"`
}

// Config is for the configuration of the database and other settings
type Config struct {
	DBURL                 string `json:"db_url"`
	DBName                string `json:"db_name"`
	TrackDBCollectionName string `json:"track_db_collection_name"`
}
