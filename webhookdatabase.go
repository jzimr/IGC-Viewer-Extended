package main

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var webhookGlobalDB WebhookMongoDB

// WebhookMongoDB stores the information of the DB connection
type WebhookMongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

/*
Init initializes the mongo storage
*/
func (db *WebhookMongoDB) Init() WebhookMongoDB {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"webhook_id"},
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
func (db *WebhookMongoDB) Add(t WebhookData) {
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
Update updates the data in the database
*/
func (db *WebhookMongoDB) Update(t WebhookData) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Update(bson.M{"id": t.ID}, t)

	if err != nil {
		fmt.Printf("error in Update(): %v", err.Error())
	}
}

/*
Delete deletes the webhook in the database
*/
func (db *WebhookMongoDB) Delete(t WebhookData) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Remove(bson.M{"id": t.ID})

	if err != nil {
		fmt.Printf("error in Delete(): %v", err.Error())
		return false
	}
	// Successfully removed webhook from DB
	return true
}

/*
Count returns the number of items currently in our database
*/
func (db *WebhookMongoDB) Count() int {
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
func (db *WebhookMongoDB) Get(keyID string) (WebhookData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	webhook := WebhookData{}

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"id": keyID}).One(&webhook)
	if err != nil {
		return webhook, false
	}
	return webhook, true
}

/*
GetAll returns all items in a collection
*/
func (db *WebhookMongoDB) GetAll() []WebhookData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allWebhooks []WebhookData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&allWebhooks)
	if err != nil {
		return []WebhookData{}
	}
	return allWebhooks
}

/*
GetLatest returns the latest added webhook in the collection
*/
func (db *WebhookMongoDB) GetLatest() WebhookData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allHooks := db.GetAll()

	if len(allHooks) == 0 {
		return WebhookData{}
	}
	return allHooks[len(allHooks)-1]
}
