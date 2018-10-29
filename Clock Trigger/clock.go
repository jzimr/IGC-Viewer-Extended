package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var webhook string
var config Config

// itIsTime checks if 10 minutes has passed
func itIsTime(timer time.Time) bool {
	if time.Now().Unix()-timer.Unix() >= 600 {
		return true
	}
	return false
}

// newTracksAdded returns all tracks that have been added since "timestamp"
// returns empty array if no tracks have been added
func changedTracksSinceLastCheck(trackCount int) int {
	allTracks := len(getAllTracks())

	difference := allTracks - trackCount
	if difference < 0 {
		difference = -difference
	}

	return difference
}

func getAllTracks() []string {
	var allTracks TracksResponse

	resp, err := http.Get("https://igcviewer-extended.herokuapp.com/paragliding/api/track")
	if err != nil {
		fmt.Println("Error when trying to make the GET request, " + err.Error())
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&allTracks)

	return allTracks
}

//////////////////////////////////////////////
// Webhook
//////////////////////////////////////////////

//	What: the ID assigned to the track that was registered
//	Response type: application/json
func invokeWebhook(changedTracks int) {
	fmt.Println("Tracks have been added or removed!")

	var dHook PostDiscordWebhook

	resp, err := http.Get("https://igcviewer-extended.herokuapp.com/paragliding/api/ticker/latest")

	if err != nil {
		fmt.Println("Error when trying to make the GET request, " + err.Error())
		return
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Some error occured, " + err.Error())
		return
	}

	timestamp := strings.TrimSuffix(string(responseData), "\n")

	// Form our Content
	dHook.Content = "Latest timestamp: " + timestamp +
		". " + strconv.Itoa(changedTracks) + " tracks have been modified. "

	endTimer := time.Since(processingTime).Nanoseconds() / int64(time.Millisecond)
	dHook.Content += "Processing time: " + strconv.FormatInt(endTimer, 10) + "ms"

	// Prepare POST request to webhook and send
	b, err := json.Marshal(dHook)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = http.Post(webhook, "application/json", strings.NewReader(string(b)))

	if err != nil {
		fmt.Println("Some error occured, " + err.Error())
	}
}

//////////////////////////////////////////////
// Configure settings (E.g. webhook URL)
//////////////////////////////////////////////

func configure() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	webhook = config.WebhookURL
}
