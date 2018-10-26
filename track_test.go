package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	igc "github.com/marni/goigc"
)

func Test_errorRubbishRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(generalHandler))
	defer server.Close()

	testUrls := []string{
		server.URL + "/paragliding/",
		server.URL + "/paragliding/rubbish",
		server.URL + "/rubbish",
		server.URL + "/paragliding/api/rubbish",
	}

	tryJSONGet(testUrls[0], http.StatusNotFound, t)
	tryJSONGet(testUrls[1], http.StatusOK, t)

	for _, testURL := range testUrls {
		resp, err := http.Get(testURL)
		if err != nil {
			t.Errorf("Error when trying to make the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected StatusCode %d, received %d", http.StatusNotFound, resp.StatusCode)
		}
	}
}

func Test_errorInvalidBodyPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	// Test with empty body
	tryJSONPost(server.URL+"/igcinfo/api/igc", "", http.StatusBadRequest, t)

	// Test for malformed url with valid JSON body
	igcJSON := "{ \"url\": \"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medell\"}"
	tryJSONPost(server.URL+"/igcinfo/api/igc", igcJSON, http.StatusBadRequest, t)
}

func Test_emptyArrayReturned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	// Check that GET /api/igc returns an empty array if no tracks stored
	resp := tryJSONGet(server.URL+"/igcinfo/api/igc", http.StatusOK, t)

	// Try to put the ids returned into an array
	var testIgcMap []string
	err := json.NewDecoder(resp.Body).Decode(&testIgcMap)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
	}

	// Length of array should be 0
	if len(testIgcMap) != 0 {
		t.Errorf("Expected an array length of 0, got %d", len(testIgcMap))
	}
}

func Test_successAddIgcFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	// Test for valid url with valid JSON body
	igcJSON := "{ \"url\": \"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc\"}"
	resp := tryJSONPost(server.URL+"/igcinfo/api/igc", igcJSON, http.StatusOK, t)

	// Check if assigned to "id0"
	var responseID RespondWithID
	err := json.NewDecoder(resp.Body).Decode(&responseID)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
		return
	}
	if responseID.ID != "id0" {
		t.Errorf("Expected id0, received: %s", responseID.ID)
	}

	// Check that GET /api/igc returns all ids in an array
	resp = tryJSONGet(server.URL+"/igcinfo/api/igc", http.StatusOK, t)

	var testIgcMap []string
	err = json.NewDecoder(resp.Body).Decode(&testIgcMap)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
	}

	if len(testIgcMap) != 1 {
		t.Errorf("Expected an array length of 1, got %d", len(testIgcMap))
	}
	// resetTest() // This is necessary or else other tests will use the data produced here.
}

//	Trying to add 3 more igc files with the same method as last test
func Test_successAddMultipleIgcFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	igcJSON := "{ \"url\": \"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc\"}"

	for i := 0; i < 3; i++ {

		//	Post the igc file to the server
		resp := tryJSONPost(server.URL+"/igcinfo/api/igc", igcJSON, http.StatusOK, t)

		// Decode the id received
		var responseID RespondWithID
		err := json.NewDecoder(resp.Body).Decode(&responseID)
		if err != nil {
			t.Errorf("Could not decode JSON returned, " + err.Error())
			return
		}

		// Check if correct ID was assigned to the track
		expectedID := "id" + strconv.Itoa(i)
		if responseID.ID != expectedID {
			t.Errorf("Expected %s, received: %s", expectedID, responseID.ID)
		}
	}

	// Check that GET /api/igc returns all ids in an array
	resp := tryJSONGet(server.URL+"/igcinfo/api/igc", http.StatusOK, t)

	// Try to put the ids returned into an array
	var testIgcMap []string
	err := json.NewDecoder(resp.Body).Decode(&testIgcMap)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
	}

	// Length of array should be 3
	if len(testIgcMap) != 3 {
		t.Errorf("Expected an array length of 3, got %d", len(testIgcMap))
	}
	// resetTest() // This is necessary or else other tests will use the data produced here.
}

// Try to request a single track
func Test_TrackRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	// Check if error returned when <id> does not exist
	resp := tryJSONGet(server.URL+"/igcinfo/api/igc/id0", http.StatusNotFound, t)

	//	Add track
	igcJSON := "{ \"url\": \"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc\"}"
	resp = tryJSONPost(server.URL+"/igcinfo/api/igc", igcJSON, http.StatusOK, t)

	// Test if we can get the track with id0
	resp = tryJSONGet(server.URL+"/igcinfo/api/igc/id0", http.StatusOK, t)

	var respTrack TrackMetaData
	var newTrack TrackMetaData

	// Decode the track we got
	err := json.NewDecoder(resp.Body).Decode(&respTrack)
	if err != nil {
		t.Errorf("Could not decode JSON returned, " + err.Error())
		return
	}

	// Trying to paste so we can compare
	track, err := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")
	if err != nil { //	If parsing failed
		t.Error("Error trying to parse a track, " + err.Error())
		return
	}
	newTrack = createTrack(track, "http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")

	// Compare results
	if respTrack.Glider != newTrack.Glider || respTrack.GliderID != newTrack.GliderID ||
		respTrack.Hdate != newTrack.Hdate || respTrack.Pilot != newTrack.Pilot ||
		respTrack.TrackLength != newTrack.TrackLength {
		t.Errorf("The results in the response are different than the actual track object, " + err.Error())
	}
	// resetTest()
}

// For GET /api/igc<id>/<field> as well
func Test_TrackFieldRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(trackHandler))
	defer server.Close()

	// Add a track
	igcJSON := "{ \"url\": \"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc\"}"
	tryJSONPost(server.URL+"/igcinfo/api/igc", igcJSON, http.StatusOK, t)

	// Check for error if <field> does not exist
	tryJSONGet(server.URL+"/igcinfo/api/igc/id0/distance", http.StatusNotFound, t)

	testFields := []string{
		"pilot",
		"glider",
		"glider_id",
		"track_length",
		"H_date",
	}

	// Check if all field requests return 200 OK and the correct data
	for _, field := range testFields {
		tryJSONGet(server.URL+"/igcinfo/api/igc/id0/"+field, http.StatusOK, t)
	}
	// resetTest()
}
