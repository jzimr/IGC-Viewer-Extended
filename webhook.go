package main

import (
	"fmt"
	"net/http"
	"strings"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")            //	parts[1]="paragliding", parts[2]="api"
	w.Header().Set("Content-Type", "application/json") // Default header type

	switch r.Method {
	case "POST":
		if len(parts) >= 4 {
			if len(parts) == 5 && parts[4] == "new_track" {

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
