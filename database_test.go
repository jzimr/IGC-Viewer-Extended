package main

import (
	"testing"

	"github.com/marni/goigc"

	"gopkg.in/mgo.v2"
)

func initDB(t *testing.T) *TrackMongoDB {
	db := TrackMongoDB{
		"mongodb://igc-app:application123@ds237373.mlab.com:37373/igcviewer-extended-db",
		"igcviewer-extended-db",
		"tracks",
	}

	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()

	if err != nil {
		t.Error(err)
	}
	return &db
}

func TestTrackMongoDB_Add(t *testing.T) {
	db := initDB(t)

	db.Init()

	tr, err := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")

	if err != nil {
		t.Error(err)
	}

	track := createTrack(tr, "http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")

	db.Add(track)

	// if db.Count() < 1 {
	// 	t.Error("adding a new track failed")
	// }

}
