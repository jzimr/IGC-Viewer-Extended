package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getDBRecords(w http.ResponseWriter) {
	fmt.Fprintln(w, trackGlobalDB.Count())
}

func deleteDBRecords(w http.ResponseWriter) {
	tracksDeleted := trackGlobalDB.Count()
	err := trackGlobalDB.DeleteAllTracks()

	if err != nil {
		fmt.Fprintln(w, "Error when trying to delete tracks, "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, tracksDeleted)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")      // parts[2]="api"
	w.Header().Set("Content-Type", "text/plain") // Default header type

	switch r.Method {
	case "GET":
		if len(parts) >= 4 && parts[2] == "api" && parts[3] == "tracks_count" {
			getDBRecords(w)
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}
	case "DELETE":
		if len(parts) >= 4 && parts[2] == "api" && parts[3] == "tracks" {
			deleteDBRecords(w)
		} else {
			http.Error(w, "Not a valid URL", http.StatusNotFound)
			return
		}
	}
}
