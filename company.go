package main

import (
	"context"
	"errors"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Company struct {
	length     int
	collection *mongo.Collection
}

func (c *Company) Add(d *DepartmentID) error {
	_, err1 := c.GetID(d.DepartmentID)
	if err1 != nil && err1.Error() != "mongo: no documents in result" {
		return errors.New("this element already exists: err:" + err1.Error())
	}
	c.DeleteID(d.DepartmentID)
	c.length++
	_, err := c.collection.InsertOne(context.TODO(), *d)
	if err != nil {
		log.Fatal("Error during insert employee in database")
	}
	return nil
}

func (c *Company) GetID(id int) (DepartmentID, error) {
	filter := bson.D{{"department_id", id}}
	// create a value into which the result can be decoded
	var result DepartmentID
	err := c.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c *Company) GetALL(db *Database) ([]Department, error) {
	cursor, err := c.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return make([]Department, 0), err
	}
	results := make([]DepartmentID, 0)
	for cursor.Next(context.TODO()) {
		var elem DepartmentID
		err := cursor.Decode(&elem)
		if err != nil {
			return make([]Department, 0), err
		}
		results = append(results, elem)
	}
	if len(results) == 0 {
		return make([]Department, 0), errors.New("there aren't records in db" + strconv.Itoa(len(results)))
	} else {
		answ := make([]Department, len(results))
		for i := 0; i < len(results); i++ {
			val, err := FillDepartment(&results[i], db)
			if err != nil {
				return answ, err
			}
			answ[i] = val
		}
		return answ, nil
	}
}

func (c *Company) DeleteID(id int) error {
	filter := bson.D{{"department_id", id}}
	num, err := c.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	if num.DeletedCount == 0 {
		return errors.New("there isn`t this element in db")
	}
	c.length--
	return nil
}

func (c *Company) DeleteALL() error {
	c.length = 0
	_, err := c.collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Company) Put(id int, d *DepartmentID) {
	filter := bson.D{{"id", id}}
	_, err := c.collection.DeleteMany(context.TODO(), filter)
	c.Add(d)
	if err != nil {
		log.Fatal(err)
	}
}
