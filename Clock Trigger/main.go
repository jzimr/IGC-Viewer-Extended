package main

import (
	"time"
)

var interval = 10

func main() {
	// Run configuration (Initialize database connection)
	configure()
	trackCount := globalDB.Count() - 1

	timer := time.Unix(time.Now().Unix()-700, 0)

	for {
		time.Sleep(10 * time.Second)

		// Has 10 minutes passed?
		if itIsTime(timer) {
			newTracks := newTracksSinceLastCheck(trackCount)

			// Do we have any new tracks?
			if len(newTracks) != 0 {
				// Send new tracks to slack webhook
				invokeWebhook(newTracks)
			}

			// reset
			timer = time.Now()
			trackCount = globalDB.Count()
		}
	}
}
