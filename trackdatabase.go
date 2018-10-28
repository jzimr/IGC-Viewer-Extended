package main

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var trackGlobalDB TrackMongoDB

// TrackMongoDB stores the information of the DB connection
type TrackMongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

/*
Init initializes the mongo storage
*/
func (db *TrackMongoDB) Init() TrackMongoDB {
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

	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return *db
}

/*
Add adds a new item to the storage
*/
func (db *TrackMongoDB) Add(t TrackMetaData) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(t)
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
	}
}

/*
DeleteAllTracks deletes all tracks in the database
*/
func (db *TrackMongoDB) DeleteAllTracks() error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	_, err = session.DB(db.DatabaseName).C(db.CollectionName).RemoveAll(nil)

	return err
}

/*
Count returns the number of items currently in our database
*/
func (db *TrackMongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}
	return count
}

/*
Get returns an item with a given ID or empty item struct
*/
func (db *TrackMongoDB) Get(keyID string) (TrackMetaData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	track := TrackMetaData{}

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"id": keyID}).One(&track)
	if err != nil {
		return track, false
	}
	return track, true
}

/*
GetAll returns all items in a collection
*/
func (db *TrackMongoDB) GetAll() []TrackMetaData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allTracks []TrackMetaData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&allTracks)
	if err != nil {
		return []TrackMetaData{}
	}
	return allTracks
}

/*
GetLatest returns the latest added trak in the collection
*/
func (db *TrackMongoDB) GetLatest() TrackMetaData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allTracks := db.GetAll()

	if len(allTracks) == 0 {
		return TrackMetaData{}
	}
	return allTracks[len(allTracks)-1]
}
