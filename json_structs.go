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

// TrackMetaInfo includes meta information about a particular track
type TrackMetaInfo struct {
	Hdate       string  `json:"H_date"`
	Pilot       string  `json:"pilot"`
	Glider      string  `json:"glider"`
	GliderID    string  `json:"glider_id"`
	TrackLength float64 `json:"track_length"` //	Does not work
}

// Empty is used to return an empty json body
type Empty struct {
}
