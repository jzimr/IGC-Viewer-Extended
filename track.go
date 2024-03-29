package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	igc "github.com/marni/goigc"
)

//	What: track registration
//	Response type: application/json
func postNewTrack(w http.ResponseWriter, r *http.Request) {
	var link IgcURL
	// Decode the link
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, "Invalid json, "+err.Error(), http.StatusBadRequest)
		return
	}
	//	Check if the URL is pointing to an IGC file
	if !strings.Contains(link.URL, ".igc") {
		http.Error(w, "Not a valid .igc file", http.StatusBadRequest)
	}

	// Convert link -> track
	track, err := igc.ParseLocation(link.URL)
	if err != nil { //	If parsing failed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add Track to database
	newTrack := createTrack(track, link.URL)
	trackGlobalDB.Add(newTrack)

	// Everything went gucci, so we reply with the ID in JSON format
	ID := RespondWithID{newTrack.ID}
	err = json.NewEncoder(w).Encode(ID)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}

	// Invoke webhooks
	invokeWebhook(w)
}

//	What: returns the array of all track ids
//	Response type: application/json
func replyWithArray(w http.ResponseWriter) {
	IDs := make([]string, 0, trackGlobalDB.Count()) // Create a new array of strings
	for i := 0; i < trackGlobalDB.Count(); i++ {
		IDs = append(IDs, "id"+strconv.Itoa(i))
	}
	err := json.NewEncoder(w).Encode(IDs)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}
}

//	What: returns the meta information about a given track with the provided <id>,
//		  or NOT FOUND response code with an empty body
//	Response type: application/json
func replyWithTrack(w http.ResponseWriter, ID string) {
	// Check if ID == ok
	track, ok := trackGlobalDB.Get(ID)
	if !ok {
		http.Error(w, "The particular ID was not found", http.StatusNotFound)
		return
	}
	metaDataView := TrackMetaDataView{track.Hdate, track.Pilot, track.Glider, track.GliderID, track.TrackLength, track.TrackSrcURL}
	// Create JSON of track meta info and return
	err := json.NewEncoder(w).Encode(metaDataView)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}
}

//	What: returns the single detailed meta information about a given track with the provided <id>,
//		  or NOT FOUND response code with an empty body. The response should always be a string,
//		  with the exception of the calculated track length, that should be a number.
//	Response type: text/plain
func replyWithTrackField(w http.ResponseWriter, ID string, field string) {
	w.Header().Set("Content-Type", "text/plain") // The response type is text/plain so we set it as this
	track, ok := trackGlobalDB.Get(ID)           // Try to get the ID requested by user
	var value string

	if !ok { // If ID was not found
		http.Error(w, "The particular ID was not found", http.StatusNotFound)
		return
	}
	var err error

	// Find the field and set the value of "value"
	switch field {
	case "pilot":
		value = track.Pilot
	case "glider":
		value = track.Glider
	case "glider_id":
		value = track.GliderID
	case "track_length":
		_, err = fmt.Fprintln(w, track.TrackLength)
	case "H_date":
		value = track.Hdate
	case "track_src_url":
		value = track.TrackSrcURL
	default:
		http.Error(w, "Not a valid <field> in the URL", http.StatusNotFound)
		return
	}
	if field != "track_length" {
		_, err = fmt.Fprintln(w, value)
	}

	if err != nil {
		http.Error(w, "Could not print data, "+err.Error(), http.StatusInternalServerError)
	}
}

// Handles POST and GET requests for /api/track/...
func trackHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") //	parts[1]="paragliding", parts[2]="api"

	switch r.Method {
	case "POST":
		if len(parts) == 4 && parts[3] == "track" { // POST /api/track
			postNewTrack(w, r)
		}

	case "GET":
		if len(parts) == 4 {
			if parts[3] == "track" { // GET /api/track
				replyWithArray(w)
			} else if parts[3] == "" { // GET /api
				replyWithServiceInfo(w)
			}
		} else if len(parts) > 4 { // /api/track...
			if len(parts) == 5 && parts[4] != "" { // GET api/track/<id>
				replyWithTrack(w, parts[4])
			} else if len(parts) == 6 && parts[5] != "" { // GET api/track/<id>/<field>
				replyWithTrackField(w, parts[4], parts[5])
			} else {
				http.Error(w, "Not a valid URL", http.StatusNotFound)
				return
			}
		}

	default:
		http.Error(w, "Not a valid request", http.StatusNotImplemented)
		return
	}
}
