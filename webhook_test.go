package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_successNewWebhook(t *testing.T) {
	resetTest()
	server := httptest.NewServer(http.HandlerFunc(webhookHandler))
	defer server.Close()

	// POST webhook to API (Omitting the "minTriggerValue")
	hookJSON := "{\"webhookURL\":\"https://discordapp.com/api/webhooks/505722994237374466/6yqNRGY1b8jitN_jyhHxLhGc-xThQBqW3L0-xC-X86Hrd__Zi_eMAGki87lv5xzbY2IQ\"}"
	resp := tryJSONPost(server.URL+"/paragliding/api/webhook/new_track", hookJSON, http.StatusCreated, t)

	// Convert response to string
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	responseString := strings.TrimSuffix(string(responseData), "\n")

	if responseString != "id0" {
		t.Errorf("Returned ID is not as expected. Expected \"id0\", got " + responseString)
	}
}

func Test_failNewWebhook(t *testing.T) {
	resetTest()
	server := httptest.NewServer(http.HandlerFunc(webhookHandler))
	defer server.Close()

	// POST some rubbish JSON and expect StatusCode 400 back
	hookJSON := "{\"URLasdasdRUBBISH>><<(!/(%()/&¤L\":\"123\"}"
	tryJSONPost(server.URL+"/paragliding/api/webhook/new_track", hookJSON, http.StatusBadRequest, t)
}

func Test_successGetWebhook(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(webhookHandler))
	defer server.Close()

	// Try to GET webhook from previous test
	resp := tryJSONGet(server.URL+"/paragliding/api/webhook/new_track/id0", http.StatusOK, t)

	// Convert response to WebhookRegistration struct
	var webhook WebhookRegistration
	json.NewDecoder(resp.Body).Decode(&webhook)

	if webhook.URL != "https://discordapp.com/api/webhooks/505722994237374466/6yqNRGY1b8jitN_jyhHxLhGc-xThQBqW3L0-xC-X86Hrd__Zi_eMAGki87lv5xzbY2IQ" &&
		webhook.MinTriggerValue != 1 {
		t.Errorf("The values returned are not identical to the JSON parsed into the API")
	}
}

func Test_failGetWebhook(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(webhookHandler))
	defer server.Close()

	// Get an id that does not exist (id5000)
	tryJSONGet(server.URL+"/paragliding/api/webhook/new_track/id5000", http.StatusNotFound, t)
}

func Test_successDeleteWebhook(t *testing.T) {
	configure()
	server := httptest.NewServer(http.HandlerFunc(webhookHandler))
	defer server.Close()

	client := server.Client()

	// Create request
	req, err := http.NewRequest("DELETE", server.URL+"/paragliding/api/webhook/new_track/id0", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to do a DELETE request, " + err.Error())
		return
	}
	defer resp.Body.Close()

	// If status is not OK
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, body: %s", http.StatusOK, resp.StatusCode, all)
	}

	// Verify deletion by trying to get "id0"
	tryJSONGet(server.URL+"/paragliding/api/webhook/new_track/id0", http.StatusNotFound, t)
}
