package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//	Helper functions

////////////////////////////////////////////////////////////////
/// Clears all data in the database
////////////////////////////////////////////////////////////////
func resetTest() {
	configure()

	err := trackGlobalDB.DeleteAllTracks()
	if err != nil {
		panic(err)
	}
	err = webhookGlobalDB.DeleteAll()
	if err != nil {
		panic(err)
	}
}

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
