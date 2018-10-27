package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//	What: Registration of new webhook for notifications about tracks being added to the system.
//	Response type: application/json
func registerNewWebhook(w http.ResponseWriter, r *http.Request) {
	var hook WebhookRegistration

	err := json.NewDecoder(r.Body).Decode(&hook)

	fmt.Println("|" + hook.URL)

	if err != nil {
		http.Error(w, "Invalid json, "+err.Error(), http.StatusBadRequest)
		return
	}

	// If ommited from POST, default value = 1
	if hook.MinTriggerValue == 0 {
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
	fmt.Fprintln(w, newhook.ID)

}

//	What: the ID assigned to the track that was registered
//	Response type: application/json
func replyWithID1(w http.ResponseWriter) {
	ID := RespondWithID{getLastID()}
	json.NewEncoder(w).Encode(ID)
}

//	What: the ID assigned to the track that was registered
//	Response type: application/json
func replyWithID2(w http.ResponseWriter) {
	ID := RespondWithID{getLastID()}
	json.NewEncoder(w).Encode(ID)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") //	parts[1]="paragliding", parts[2]="api"

	switch r.Method {
	case "POST":
		if len(parts) >= 5 {
			if parts[4] == "new_track" { // Fix on rubbish req
				registerNewWebhook(w, r)
			}
		}

	case "GET":
		if len(parts) >= 4 && len(parts) < 6 {
			if len(parts) == 4 || parts[4] == "" {
				replyWithTicker(w)
			} else if len(parts) == 5 && parts[4] == "latest" {
				replyWithLatest(w)
			} else if len(parts) == 5 {
				fmt.Println() // Placeholder // replyWithTimestamp(w)
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
