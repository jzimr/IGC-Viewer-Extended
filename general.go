package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

//	What: meta information about the API
//	Response type: application/json
func replyWithServiceInfo(w http.ResponseWriter) {
	metaInfo := MetaInfo{Uptime: getUptime(), Info: "Service for IGC tracks.", Version: "v1"}

	err := json.NewEncoder(w).Encode(metaInfo)
	if err != nil {
		http.Error(w, "Could not encode json, "+err.Error(), http.StatusInternalServerError)
	}
}

// Handles GET request for /api
func generalHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") //	parts[1]="paragliding", parts[2]="api"

	switch r.Method {
	case "GET":
		if (len(parts) == 4 && parts[2] == "api" && parts[3] == "") ||
			(len(parts) == 3 && parts[2] == "") {
			replyWithServiceInfo(w)
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}

	default:
		http.Error(w, "Not a valid request", http.StatusNotImplemented)
		return
	}
}
