package main

import (
	"strconv"
	"time"

	igc "github.com/marni/goigc"
)

////////////////////////////////////////////////////////////////
/// Gets the uptime of the global "uptime" variable
////////////////////////////////////////////////////////////////
func getUptime() string {
	uYear, uMonth, uDay := uptime.Date()     // Get stored Y, M, D
	nYear, nMonth, nDay := time.Now().Date() // Get current Y, M, D
	interval := time.Since(uptime)           // Get H, M, S

	years := strconv.Itoa(nYear - uYear)
	months := strconv.Itoa(int(nMonth) - int(uMonth))
	days := strconv.Itoa(nDay - uDay)
	hours := strconv.Itoa(int(interval.Hours()) % 24)
	minutes := strconv.Itoa(int(interval.Minutes()) % 60)
	seconds := strconv.Itoa(int(interval.Seconds()) % 60)

	return "P" + years + "Y" + months + "M" + days + "DT" + hours + "H" + minutes + "M" + seconds + "S"
}

////////////////////////////////////////////////////////////////
//	Function that converts a Track object to a TrackMetaData object
////////////////////////////////////////////////////////////////
func createTrack(track igc.Track, URL string) TrackMetaData {
	var trackData TrackMetaData

	trackDist := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		trackDist += track.Points[i].Distance(track.Points[i+1])
	}

	newID := "id" + strconv.Itoa(trackGlobalDB.Count())

	trackData = TrackMetaData{newID, track.Date.String(), track.Pilot, track.GliderType, track.GliderID, trackDist, URL, getNow()}

	return trackData
}

////////////////////////////////////////////////////////////////
//	Function that converts a WebhookRegistration object to a
//	WebhookData object
////////////////////////////////////////////////////////////////
func createWebhook(webhook WebhookRegistration) WebhookData {
	var hookData WebhookData

	newID := "id" + strconv.Itoa(webhookGlobalDB.Count())
	hookData = WebhookData{newID, webhook.URL, webhook.MinTriggerValue}

	return hookData
}

////////////////////////////////////////////////////////////////
//	Function that returns the current UNIX time
////////////////////////////////////////////////////////////////
func getNow() int64 {
	return time.Now().Unix()
}

////////////////////////////////////////////////////////////////
//	Function that gets the last ID in database
////////////////////////////////////////////////////////////////
func getLastID() string {
	trackCount := trackGlobalDB.Count() - 1

	if trackCount < 0 { // If we don't have any tracks yet
		return "-1"
	}

	return "id" + strconv.Itoa(trackGlobalDB.Count()-1)
}

////////////////////////////////////////////////////////////////
//	Function that returns the timestamp of the latest track added
////////////////////////////////////////////////////////////////
func getLatestTimestamp() int64 {
	latestTrack, ok := trackGlobalDB.Get(getLastID())
	latestTrack2 := latestTrack.(TrackMetaData)

	if !ok {
		return 0
	}

	return latestTrack2.Timestamp
}

////////////////////////////////////////////////////////////////
//	Function that starts the counter
////////////////////////////////////////////////////////////////
func startMillCounter() time.Time {
	return time.Now()
}

////////////////////////////////////////////////////////////////
//	Function that ends the counter
////////////////////////////////////////////////////////////////
func stopMillCounter(timer time.Time) int64 {
	duration := time.Since(timer)
	return duration.Nanoseconds() / int64(time.Millisecond)
}
