package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	igc "github.com/marni/goigc"
)

//	Store information on server
var igcMap = make(map[string]igc.Track) // E.g. ["id0"], ["id1"], ...

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
	track, err := igc.ParseLocation(link.URL)
	if err != nil { //	If parsing failed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	igcMap["id"+strconv.Itoa(len(igcMap))] = track // Add track to storage

	// Check if item was successfully added to map
	_, ok := igcMap["id"+strconv.Itoa(len(igcMap)-1)]
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Everything went gucci, so we reply with the ID in JSON format
	replyWithID(w)
}

//	What: meta information about the API
//	Response type: application/json
func replyWithServiceInfo(w http.ResponseWriter) {
	metaInfo := MetaInfo{Uptime: getUptime(), Info: "Service for IGC tracks.", Version: "v1"}
	json.NewEncoder(w).Encode(metaInfo)
}

//	What: the ID assigned to the track that was registered
//	Response type: application/json
func replyWithID(w http.ResponseWriter) {
	ID := RespondWithID{"id" + strconv.Itoa(len(igcMap)-1)}
	json.NewEncoder(w).Encode(ID)
}

//	What: returns the array of all track ids
//	Response type: application/json
func replyWithArray(w http.ResponseWriter) {
	IDs := make([]string, 0, len(igcMap)) // Create a new array of strings
	for k := range igcMap {               // Fill array with Ids
		IDs = append(IDs, k)
	}
	json.NewEncoder(w).Encode(IDs)
}

//	What: returns the meta information about a given track with the provided <id>,
//		  or NOT FOUND response code with an empty body
//	Response type: application/json
func replyWithTrack(w http.ResponseWriter, ID string) {
	track, ok := igcMap[ID]
	var trackInfo TrackMetaInfo

	if !ok {
		http.Error(w, "The particular ID was not found", http.StatusNotFound)
		return
	}

	trackInfo = createMetaTrack(track)

	json.NewEncoder(w).Encode(trackInfo)
}

//	What: returns the single detailed meta information about a given track with the provided <id>,
//		  or NOT FOUND response code with an empty body. The response should always be a string,
//		  with the exception of the calculated track length, that should be a number.
//	Response type: text/plain
func replyWithTrackField(w http.ResponseWriter, ID string, field string) {
	w.Header().Set("Content-Type", "text/plain") // The response type is text/plain so we set it as this
	track, ok := igcMap[ID]                      // Try to get the ID requested by user
	var value string

	if !ok { // If ID was not found
		http.Error(w, "The particular ID was not found", http.StatusNotFound)
		json.NewEncoder(w).Encode(Empty{})
		return
	}

	// Find the field and set the value of "value"
	switch field {
	case "pilot":
		value = track.Pilot
		break
	case "glider":
		value = track.GliderType
		break
	case "glider_id":
		value = track.GliderID
		break
	case "track_length":
		break
	case "H_date":
		value = track.Date.String()
		break
	default:
		http.Error(w, "Not a valid <field> in the URL", http.StatusNotFound)
		return
	}

	if field == "track_length" {
		trackDist := 0.0
		for i := 0; i < len(track.Points)-1; i++ {
			trackDist += track.Points[i].Distance(track.Points[i+1])
		}
		fmt.Fprintln(w, trackDist)
	} else {
		fmt.Fprintln(w, value)
	}
}

// Handles all POST and GET requests from /igcinfo/api
func handlerIgc(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")            //	parts[1]="igcinfo", parts[2]="api"
	w.Header().Set("Content-Type", "application/json") // Default header type

	switch r.Method {
	case "POST":
		if len(parts) == 4 && parts[3] == "igc" { // POST /api/igc
			postNewTrack(w, r)
		}

	case "GET":
		if len(parts) == 4 && parts[3] == "" { // GET /api
			replyWithServiceInfo(w)
		} else if len(parts) >= 4 && parts[3] == "igc" { // /api/igc...
			if len(parts) == 4 { // GET api/igc
				replyWithArray(w)
			} else if len(parts) == 5 && parts[4] != "" { // GET api/igc/<id>
				replyWithTrack(w, parts[4])
			} else if len(parts) == 6 && parts[4] != "" && parts[5] != "" { // GET api/igc/<id>/<field>
				replyWithTrackField(w, parts[4], parts[5])
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
