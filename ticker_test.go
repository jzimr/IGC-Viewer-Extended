package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	igc "github.com/marni/goigc"
)

////////////////////////////////////////////////////////////////
/// Populates the track database
/// !!! Run this before doing all other tests in ticker !!!
////////////////////////////////////////////////////////////////
func Test_populateDB(t *testing.T) {
	resetTest()
	tracks := []string{
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc",
	}

	// Start DB connection
	configure()

	// Add each track to our database (TOTAL: 3)
	for _, tURL := range tracks {
		track, err := igc.ParseLocation(tURL)
		if err != nil { //	If parsing failed
			panic(err)
		}

		// Add Track to database
		newTrack := createTrack(track, tURL)
		trackGlobalDB.Add(newTrack)
	}
}

func Test_successGetLatestTicker(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer server.Close()

	// Getting timestamp from our API
	resp := tryJSONGet(server.URL+"/paragliding/api/ticker/latest", http.StatusOK, t)
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	responseString := strings.TrimSuffix(string(responseData), "\n")
	APITimestamp, err := strconv.ParseInt(responseString, 0, 64)
	if err != nil {
		t.Error(err)
	}

	// Get latest timestamp from database
	latestTimestamp := trackGlobalDB.GetLatest().Timestamp

	if APITimestamp != latestTimestamp {
		t.Error("Timestamps do not match between the database and our API!")
	}
}

func Test_successGetTracks(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer server.Close()

	// Get API and decode json
	resp := tryJSONGet(server.URL+"/paragliding/api/ticker", http.StatusOK, t)
	var APIInfo TickerInfo
	DBInfo := trackGlobalDB.GetLatest()

	err := json.NewDecoder(resp.Body).Decode(&APIInfo)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
	}

	// Check that t_latest is latest in database, and t_stop is newest track
	if DBInfo.Timestamp != APIInfo.Latest || DBInfo.Timestamp != APIInfo.Stop {
		t.Errorf("Timestamps do not match!")
	}

	// Check that t_start is oldest track
	DBInfo, ok := trackGlobalDB.Get("id0")
	if !ok {
		t.Errorf("Could not find track with id0")
	}

	if DBInfo.Timestamp != APIInfo.Start {
		t.Errorf("t_stop timestamps do not match!")
	}

	if len(APIInfo.Tracks) != 3 {
		t.Errorf("Expected 3 tracks in list, got " + strconv.Itoa(len(APIInfo.Tracks)))
	}
}

func Test_successGetTracksHigherThan(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer server.Close()

	// Get the second track (out of three) in list
	DBInfo, ok := trackGlobalDB.Get("id1")
	if !ok {
		t.Errorf("Could not find track with id1")
	}
	minTimestamp := strconv.FormatInt(DBInfo.Timestamp-1, 10)

	// Get API and decode json
	resp := tryJSONGet(server.URL+"/paragliding/api/ticker/"+minTimestamp, http.StatusOK, t)
	var APIInfo TickerInfo

	err := json.NewDecoder(resp.Body).Decode(&APIInfo)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
	}

	// Check that t_start is oldest track in the list HIGHER THAN the timestamp provided
	if DBInfo.Timestamp != APIInfo.Start {
		t.Errorf("t_stop timestamps do not match!" + strconv.FormatInt(DBInfo.Timestamp, 10) + " | " + strconv.FormatInt(APIInfo.Start, 10))
	}
}
