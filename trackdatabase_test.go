package main

import (
	"testing"

	"github.com/marni/goigc"
)

func setupDB(t *testing.T) *TrackMongoDB {
	db := TrackMongoDB{
		"mongodb://igc-app:application123@ds237373.mlab.com:37373/igcviewer-extended-db",
		"igcviewer-extended-db",
		"tracks",
	}

	db = db.Init()
	return &db
}

func Test_TrackMongoDBAdd(t *testing.T) {
	db := setupDB(t)
	db.DeleteAllTracks()

	tr, err := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")
	if err != nil {
		t.Error(err)
	}

	track := createTrack(tr, "http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")

	db.Add(track)

	if db.Count() < 1 {
		t.Error("adding a new track failed")
	}
}

func Test_TrackMongoDBDeleteAll(t *testing.T) {
	db := setupDB(t)

	// We need something to delete
	if db.Count() == 0 {
		Test_TrackMongoDBAdd(t)
		Test_TrackMongoDBAdd(t)
	}

	db.DeleteAllTracks()

	if db.Count() != 0 {
		t.Error("removing all tracks failed")
	}
}

func Test_TrackMongoDBGet(t *testing.T) {
	db := setupDB(t)

	// We need something to delete
	if db.Count() == 0 {
		Test_TrackMongoDBAdd(t)
	}

	track, ok := db.Get("id0")

	if !ok {
		t.Error("Something went wrong when trying to retrieve track id0")
	}
	if track.ID != "id0" {
		t.Error("We got a track, but not the one we wanted :(")
	}
}
