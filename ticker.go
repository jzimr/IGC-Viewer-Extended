package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//	What: returns the timestamp of the latest added track
//	Response type: text/plain
func replyWithLatest(w http.ResponseWriter) {
	// Set header
	w.Header().Set("Content-Type", "text/plain") // Default header type

	lTimestamp := getLatestTimestamp()

	if lTimestamp == -1 {
		http.Error(w, "No latest tracks yet", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, lTimestamp)
}

//	What: returns the JSON struct representing the ticker for the IGC tracks. The first track
//		  returned should be the oldest. The array of track ids returned should be capped at 5
//		  to emulate "paging" of the responses. The acp(5) should be a configuration parameter
//		  of the application(ie. easy to change by the administrator).
//	Response type: application/json
func replyWithTicker(w http.ResponseWriter) {
	timer := startMillCounter()
	MAX := config.MaxTracksPerPage

	var startStamp int64 = -1
	var stopStamp int64 = -1

	allTracks := trackGlobalDB.GetAll()

	// Return empty if database has no data
	if len(allTracks) == 0 {
		emptyTickerInfo := TickerInfo{}
		emptyTickerInfo.Tracks = make([]string, 0, 0)
		json.NewEncoder(w).Encode(emptyTickerInfo)
		return
	}

	chosenTracks := make([]string, 0, MAX) // Create a new array of strings
	for i := 0; i < len(allTracks); i++ {
		if i >= MAX { //	We don't want more results than defined
			break
		}
		chosenTracks = append(chosenTracks, allTracks[i].ID)
	}

	// Choose the oldest and newest tracks from our list
	if len(allTracks) > 0 {
		startStamp = allTracks[0].Timestamp
	}
	if len(allTracks) >= MAX {
		stopStamp = allTracks[MAX-1].Timestamp
	} else if len(allTracks) < MAX {
		stopStamp = allTracks[len(allTracks)-1].Timestamp
	}

	tInfo := TickerInfo{getLatestTimestamp(), startStamp, stopStamp,
		chosenTracks, stopMillCounter(timer)}

	json.NewEncoder(w).Encode(tInfo)
}

// 	What: returns the JSON struct representing the ticker for the IGC tracks. The first returned
// 		  track should have the timestamp HIGHER than the one provided in the query. The array
//		  of track IDs returned should be capped at 5, to emulate "paging" of the responses.
//		  The cap (5) should be a configuration parameter of the application
//		  (ie. easy to change by the administrator).
//	Response type: application/json
func replyWithTimestamp(w http.ResponseWriter, fromStamp string) {
	timestamp, err := strconv.ParseInt(fromStamp, 10, 64)
	MAX := config.MaxTracksPerPage

	if err != nil {
		http.Error(w, "Could not convert parameter <timestamp> to number", http.StatusBadRequest)
	}

	timer := startMillCounter()

	allTracks := trackGlobalDB.GetAll()

	// Find the first index where the timestamp is lower
	var fromIndex int
	for i, track := range allTracks {
		if timestamp < track.Timestamp {
			fromIndex = i
		}
	}

	chosenTracks := make([]string, 0, MAX) // Create a new array of strings
	for i := fromIndex; i < fromIndex+MAX; i++ {
		if i < len(allTracks) {
			chosenTracks = append(chosenTracks, allTracks[i].ID)
		} else {
			break
		}
	}

	tInfo := TickerInfo{getLatestTimestamp(), allTracks[0].Timestamp, allTracks[MAX-1].Timestamp,
		chosenTracks, stopMillCounter(timer)}

	json.NewEncoder(w).Encode(tInfo)
}

func tickerHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") //	parts[1]="paragliding", parts[2]="api"

	fmt.Println(r.URL.Path)
	fmt.Println(len(parts))

	switch r.Method {
	case "GET":
		if len(parts) >= 4 && len(parts) < 6 {
			if len(parts) == 4 || parts[4] == "" {
				replyWithTicker(w)
			} else if len(parts) == 5 && parts[4] == "latest" {
				replyWithLatest(w)
			} else if len(parts) == 5 {
				replyWithTimestamp(w, parts[4])
			} else {
				http.Error(w, "Not a valid URL", http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}

	default:
		http.Error(w, "Not a valid request", http.StatusNotImplemented)
		return
	}
}
