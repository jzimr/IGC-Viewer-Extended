package main

import (
	"fmt"
	"time"
)

var interval = 10
var processingTime time.Time

func main() {
	// Run configuration (Initialize database connection)
	configure()
	trackCount := len(getAllTracks())

	fmt.Println("Running...")

	timer := time.Now()

	for {
		time.Sleep(10 * time.Second)

		// Has 10 minutes passed?
		if itIsTime(timer) {
			processingTime = time.Now() // Start our timer
			changedTracks := changedTracksSinceLastCheck(trackCount)

			// Do we have any new tracks?
			if changedTracks != 0 {
				// Send new tracks to slack webhook
				invokeWebhook(changedTracks)
			}

			// reset
			timer = time.Now()
			trackCount = len(getAllTracks())
		}
	}
}
