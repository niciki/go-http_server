package main

type Employee struct {
	Id      int    `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	Salary  int    `json:"salary" bson:"salary"`
	Age     int    `json:"age" bson:"age"`
}
