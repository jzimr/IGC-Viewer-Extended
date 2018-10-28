package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var uptime time.Time //	Keeps track of how long our server has been running

// Get the port of the webserver (In this case from Heroku)
func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Run configuration (E.g. initialize database connection)
	configure()

	http.HandleFunc("/paragliding/", generalHandler) // general_api.go
	http.HandleFunc("/paragliding/api/", forwardingHandler)
	uptime = time.Now() // Start timer
	if err := http.ListenAndServe(addr /*":8080"*/, nil); err != nil {
		panic(err)
	}
}

// This functions is to forward the URL to the correct files and functions
func forwardingHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")            //	parts[1]="paragliding", parts[2]="api"
	w.Header().Set("Content-Type", "application/json") // Default header type

	if parts[3] == "" { // general_api.go
		generalHandler(w, r)
	} else if parts[3] == "track" { // track_api.go
		trackHandler(w, r)
	} else if parts[3] == "ticker" { // ticker_api.go
		tickerHandler(w, r)
	} else if parts[3] == "webhook" && parts[4] == "new_track" { // webhook_api.go
		webhookHandler(w, r)
	}
}
