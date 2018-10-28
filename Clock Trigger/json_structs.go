package main

// PostDiscordWebhook contains the message we want to send to the discord webhook
type PostDiscordWebhook struct {
	Content string `json:"content"`
}

// TracksResponse includes a list of track IDs
// type TracksResponse struct {
// 	Tracks []string
// }

// TracksResponse includes a list of track IDs
type TracksResponse []string

// Config is for the configuration of the database and other settings
type Config struct {
	WebhookURL string `json:"webhook_url"`
}
