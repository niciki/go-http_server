package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	length     int
	m          sync.RWMutex
	collection *mongo.Collection
}

func (db *Database) Add(e *Employee) error {
	_, err1 := db.GetID(e.Id)
	if err1 != nil && err1.Error() != "mongo: no documents in result" {
		return errors.New("this element already exists: err:" + err1.Error())
	}
	db.DeleteID(e.Id)
	db.length++
	_, err := db.collection.InsertOne(context.TODO(), *e)
	if err != nil {
		log.Fatal("Error during insert employee in database")
	}
	return nil
}

func (db *Database) GetID(id int) (Employee, error) {
	filter := bson.D{{"id", id}}
	// create a value into which the result can be decoded
	var result Employee
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *Database) GetAll() ([]Employee, error) {
	cursor, err := db.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return make([]Employee, 0), err
	}
	results := make([]Employee, 0)
	for cursor.Next(context.TODO()) {
		var elem Employee
		err := cursor.Decode(&elem)
		if err != nil {
			return results, err
		}
		results = append(results, elem)
	}
	if len(results) == 0 {
		return make([]Employee, 0), errors.New("there aren't records in db" + strconv.Itoa(len(results)))
	} else {
		return results, nil
	}
}

func (db *Database) Put(id int, e *Employee) {
	filter := bson.D{{"id", id}}
	_, err := db.collection.DeleteMany(context.TODO(), filter)
	db.Add(e)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *Database) DeleteID(id int) error {
	filter := bson.D{{"id", id}}
	num, err := db.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	if num.DeletedCount == 0 {
		return errors.New("there isn`t this element in db")
	}
	db.length--
	return nil
}

func (db *Database) DeleteALL() error {
	db.length = 0
	_, err := db.collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func NewDatabase(collection *mongo.Collection) *Database {
	db := Database{
		collection: collection,
	}
	leng, _ := db.GetAll()
	db.length = len(leng)
	return &db
}
