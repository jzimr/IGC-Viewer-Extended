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
//	Function that converts a Track object to a TrackMetaInfo object
////////////////////////////////////////////////////////////////
func createMetaTrack(track igc.Track) TrackMetaInfo {
	var trackInfo TrackMetaInfo

	trackDist := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		trackDist += track.Points[i].Distance(track.Points[i+1])
	}
	trackInfo = TrackMetaInfo{track.Date.String(), track.Pilot, track.GliderType, track.GliderID, trackDist}

	return trackInfo
}
