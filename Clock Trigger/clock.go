package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var webhook = "https://hooks.slack.com/services/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"
var config Config
var globalDB MongoDB

// itIsTime checks if 10 minutes has passed
func itIsTime(timer time.Time) bool {
	if time.Now().Unix()-timer.Unix() >= 600 {
		return true
	}
	return false
}

// newTracksAdded returns all tracks that have been added since "timestamp"
// returns empty array if no tracks have been added
func newTracksSinceLastCheck(trackCount int) []TrackMetaData {
	var newTracks []TrackMetaData

	if globalDB.Count() > trackCount {
		tracks := globalDB.GetAll()

		for i := trackCount; i < len(tracks); i++ {
			newTracks = append(newTracks, tracks[i])
		}
	}
	return newTracks
}

//////////////////////////////////////////////
// Webhook
//////////////////////////////////////////////

//	What: the ID assigned to the track that was registered
//	Response type: application/json
func invokeWebhook(tracks []TrackMetaData) {
	startTimer := time.Now()

	var sHook PostSlackWebhook
	latestTrack := tracks[len(tracks)-1]

	sHook.Text = "Latest timestamp: " + strconv.FormatInt(latestTrack.Timestamp, 10) +
		", " + strconv.Itoa(len(tracks)) + " new tracks: "

	for _, t := range tracks {
		sHook.Text += t.ID + ", "
	}

	endTimer := time.Since(startTimer).Nanoseconds() / int64(time.Millisecond)
	sHook.Text += "Processing time: " + strconv.FormatInt(endTimer, 10) + "ms"

	b, err := json.Marshal(sHook)
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
// Database
//////////////////////////////////////////////

/*
Init initializes the mongo storage
*/
func (db *MongoDB) Init() MongoDB {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"track_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return *db
}

/*
Count returns the number of items currently in our database
*/
func (db *MongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}
	return count
}

/*
GetAll returns all items in a collection
*/
func (db *MongoDB) GetAll() []TrackMetaData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allTracks []TrackMetaData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&allTracks)
	if err != nil {
		return []TrackMetaData{}
	}
	return allTracks
}

//////////////////////////////////////////////
// Configure database connection
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

	// Start a connection to our track global database
	globalDB = MongoDB{
		config.DBURL,
		config.DBName,
		config.TrackDBCollectionName,
	}
	globalDB = globalDB.Init()
}
