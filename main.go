package main

import (
	"fmt"
	"net/http"
	"os"
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
	// addr, err := determineListenAddress()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	http.HandleFunc("/igcinfo/api/", handlerIgc)
	uptime = time.Now() // Start timer
	if err := http.ListenAndServe( /*add*/ ":8080", nil); err != nil {
		panic(err)
	}
}
