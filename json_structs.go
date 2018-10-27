package main

// MetaInfo holds information about our system
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// IgcURL is the POST url of the igc url the user sent
type IgcURL struct {
	URL string `json:"url"`
}

// RespondWithID holds the ID of the Track added by the client
type RespondWithID struct {
	ID string `json:"id"`
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

// TrackMetaDataView shows a piece of the data to the client
type TrackMetaDataView struct {
	Hdate       string  `json:"H_date"`
	Pilot       string  `json:"pilot"`
	Glider      string  `json:"glider"`
	GliderID    string  `json:"glider_id"`
	TrackLength float64 `json:"track_length"`
	TrackSrcURL string  `json:"track_src_url"`
}

// TickerInfo contains all data about the latest added timestamp, first timestamp of track, ...
type TickerInfo struct {
	Latest     int64    `json:"t_latest"`
	Start      int64    `json:"t_start"`
	Stop       int64    `json:"t_stop"`
	Tracks     []string `json:"tracks"`
	Processing int64    `json:"processing"`
}

// WebhookRegistration contains the information of the new webhook the client wants to register
type WebhookRegistration struct {
	URL             string `json:"webhookURL"`
	MinTriggerValue int    `json:"minTriggerValue"`
}

// WebhookData contains the information of the new webhook the client wants to register
type WebhookData struct {
	ID              string `json:"id"`
	URL             string `json:"webhookURL"`
	MinTriggerValue int    `json:"minTriggerValue"`
}
