package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//	What: Registration of new webhook for notifications about tracks being added to the system.
//	Response type: application/json
func registerNewWebhook(w http.ResponseWriter, r *http.Request) {
	var hook WebhookRegistration

	err := json.NewDecoder(r.Body).Decode(&hook)

	if err != nil {
		http.Error(w, "Malformed JSON body, "+err.Error(), http.StatusBadRequest)
		return
	}
	if hook.URL == "" { // We don't allow empty URLs
		http.Error(w, "Malformed JSON body", http.StatusBadRequest)
	}

	// If ommited from POST, default value = 1
	if hook.MinTriggerValue < 1 {
		hook.MinTriggerValue = 1
	}

	newhook := createWebhook(hook)
	webhookGlobalDB.Add(newhook)

	fmt.Println(newhook.ID)

	// Check if item was successfully added to database
	_, ok := webhookGlobalDB.Get(newhook.ID)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Everything went gucci, so we reply with the ID
	w.Header().Set("Content-Type", "text/plain") // Set header type to text/plain
	w.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprintln(w, newhook.ID)
	if err != nil {
		http.Error(w, "Could not print value, "+err.Error(), http.StatusInternalServerError)
	}

}

//	What: Goes through all webhooks and if applicable, POST tracks to the webhook
//	Response type: application/json
func invokeWebhook(w http.ResponseWriter) {
	webhooks := webhookGlobalDB.GetAll()
	tracks := trackGlobalDB.GetAll()
	startTimer := startMillCounter()

	for _, hook := range webhooks {

		if (len(tracks)-hook.CountFromTrack)%hook.MinTriggerValue == 0 {
			var dHook PostDiscordWebhook
			latestTrack := tracks[len(tracks)-1]

			dHook.Content = "Latest timestamp: " + strconv.FormatInt(latestTrack.Timestamp, 10) +
				", " + strconv.Itoa(hook.MinTriggerValue) + " new tracks: "

			for i := len(tracks) - hook.MinTriggerValue; i < len(tracks); i++ {
				dHook.Content += "id" + strconv.Itoa(i) + ", "
			}

			endTimer := stopMillCounter(startTimer)
			dHook.Content += "Processing time: " + strconv.FormatInt(endTimer, 10) + "ms"

			b, err := json.Marshal(dHook)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = http.Post(hook.URL, "application/json", strings.NewReader(string(b)))

			if err != nil {
				http.Error(w, "Could not POST data to webhook, "+err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

//	What: Accessing registered webhooks.
//	Response type: application/json
func accessWebhook(w http.ResponseWriter, webhookID string) {
	webhook, ok := webhookGlobalDB.Get(webhookID)

	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	webhookInfo := WebhookRegistration{webhook.ID, webhook.MinTriggerValue}

	err := json.NewEncoder(w).Encode(webhookInfo)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}
}

//	What: Deleting registered webhooks.
//	Response type: application/json
func deleteWebhook(w http.ResponseWriter, webhookID string) {
	webhook, ok := webhookGlobalDB.Get(webhookID)

	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	webhookInfo := WebhookRegistration{webhook.ID, webhook.MinTriggerValue}

	ok = webhookGlobalDB.Delete(webhook)
	if !ok {
		http.Error(w, "Could not delete webhook, are you sure it exists?", http.StatusBadRequest)
	}

	err := json.NewEncoder(w).Encode(webhookInfo)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}

}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") //	parts[1]="paragliding", parts[2]="api"

	switch r.Method {
	case "POST":
		if len(parts) >= 5 && parts[4] == "new_track" {
			registerNewWebhook(w, r)
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}

	case "GET":
		if len(parts) >= 6 && parts[5] != "" {
			accessWebhook(w, parts[5])
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}
	case "DELETE":
		if len(parts) >= 6 && parts[5] != "" {
			deleteWebhook(w, parts[5])
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}

	default:
		http.Error(w, "Not a valid request", http.StatusNotImplemented)
		return
	}
}
