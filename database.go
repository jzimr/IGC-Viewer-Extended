package main

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var globalDB TrackMongoDB

// TrackMongoDB stores the information of the DB connection
type TrackMongoDB struct {
	DatabaseURL         string
	DatabaseName        string
	TrackCollectionName string
}

/*
Init initializes the mongo storage
*/
func (db *TrackMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"track_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	globalDB = *db
}

/*
Add adds a new track to the storage
*/
func (db *TrackMongoDB) Add(t TrackMetaData) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Insert(t)
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
	}
}

/*
Count returns the number of tracks currently in our database
*/
func (db *TrackMongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.TrackCollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}
	return count
}

/*
Get returns a track with a given ID or empty track struct
*/
func (db *TrackMongoDB) Get(keyID string) (TrackMetaData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	track := TrackMetaData{}

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Find(bson.M{"trackid": keyID}).One(&track)
	if err != nil {
		return track, false
	}
	return track, true
}

/*
GetAll returns all tracks
*/
func (db *TrackMongoDB) GetAll() []TrackMetaData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allTracks []TrackMetaData

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Find(bson.M{}).All(&allTracks)
	if err != nil {
		return []TrackMetaData{}
	}
	return allTracks
}
