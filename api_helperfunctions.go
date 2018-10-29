package main

import (
	"fmt"
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

	count := trackGlobalDB.Count()

	// Fixes error where trackGlobalDB.Count() sometimes returns -1 for some reason
	if count < 0 {
		count = 0
	}

	newID := "id" + strconv.Itoa(count)

	trackData = TrackMetaData{newID, track.Date.String(), track.Pilot, track.GliderType, track.GliderID, trackDist, URL, time.Now().Unix()}

	return trackData
}

////////////////////////////////////////////////////////////////
//	Function that converts a WebhookRegistration object to a
//	WebhookData object
////////////////////////////////////////////////////////////////
func createWebhook(webhook WebhookRegistration) WebhookData {
	var hookData WebhookData
	var newID string

	latestHook := webhookGlobalDB.GetLatest()

	// Assign a new ID
	// If we have no hooks yet
	if latestHook.ID == (WebhookData{}.ID) {
		newID = "id0"
	} else {
		lastID, err := strconv.Atoi(latestHook.ID[2:])
		if err != nil {
			_ = fmt.Errorf("Could not convert integer to string, " + err.Error())
		}
		newID = "id" + strconv.Itoa(lastID+1)
	}
	hookData = WebhookData{newID, webhook.URL, webhook.MinTriggerValue, trackGlobalDB.Count() - 1}

	return hookData
}

////////////////////////////////////////////////////////////////
//	Function that returns the timestamp of the latest track added
////////////////////////////////////////////////////////////////
func getLatestTimestamp() int64 {
	latestTrack := trackGlobalDB.GetLatest()

	if latestTrack.ID == "" {
		return -1
	}

	return latestTrack.Timestamp
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
