package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//	Helper functions

////////////////////////////////////////////////////////////////
/// Resets the map
////////////////////////////////////////////////////////////////
// func resetTest() {
// 	trackMap = make(map[string]igc.Track)
// 	metaTrackMap = make(map[string]TrackMetaData)
// }

////////////////////////////////////////////////////////////////
/// Tries to do a POST request to the URL and checks for errors
/// Returns: httpresponse
////////////////////////////////////////////////////////////////
// URL: 			The URL to post to
// jsonValue:		The json we want to post
// expectedStatus: 	The statuscode we expect to receive back
// t: 				The current test, so we can throw errors
/////////////////////////////////////////////////////////////
func tryJSONPost(URL string, jsonValue string, expectedStatus int, t *testing.T) (resp *http.Response) {
	resp, err := http.Post(URL, "application/json", strings.NewReader(jsonValue))
	if err != nil {
		t.Errorf("Error trying to create a POST request, %s", err)
	}
	if resp.StatusCode != expectedStatus {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, body: %s", expectedStatus, resp.StatusCode, all)
	}
	return resp
}

////////////////////////////////////////////////////////////////
/// Tries to do a GET request to the URL and checks for errors
/// Returns: httpresponse
////////////////////////////////////////////////////////////////
// URL: 			The URL to post to
// expectedStatus: 	The statuscode we expect to receive back
// t: 				The current test, so we can throw errors
/////////////////////////////////////////////////////////////
func tryJSONGet(URL string, expectedStatus int, t *testing.T) (resp *http.Response) {
	resp, err := http.Get(URL)
	if err != nil {
		t.Errorf("Error when trying to make the GET request, %s", err)
	}
	if resp.StatusCode != expectedStatus {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, body: %s", expectedStatus, resp.StatusCode, all)
	}
	return resp
}
