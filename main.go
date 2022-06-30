package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	// Get a handle for your collection
	collection_employee := client.Database("server").Collection("employee")
	collection_department := client.Database("server").Collection("department")
	handler := NewHandler(collection_employee, collection_department)
	router := gin.Default()
	// Add new employee
	router.POST("/employee", handler.PostFuncEmployee)
	// Add new department
	router.POST("/department", handler.PostFuncDepartment)
	// Get info about employee
	router.GET("/employee/:id", handler.GetFuncEmployeeID)
	// Get array with all data of database
	router.GET("/employee", handler.GetFuncEmployeeALL)
	// Get info about department by id
	router.GET("/department/:id", handler.GetFuncDepartmentID)
	// Get array with all departments
	router.GET("/department", handler.GetFuncDepartmentALL)
	// Get stat about company
	router.GET("/statistic", handler.GetStat)
	// Change facts about employee
	router.PUT("/employee/:id", handler.PutFuncEmployee)
	// Add to Department with number id employee with employeid
	router.PUT("/department/add/:id/:employeeid", handler.PutFuncAddDepartment)
	// Delete from Department with number id employee with  employeid
	router.PUT("/department/delete/:id/:employeeid", handler.PutFuncDeleteDepartment)
	// Delete an employee
	router.DELETE("/employee/:id", handler.DeleteFuncEmployeeID)
	// Delete all employees
	router.DELETE("/employee", handler.DeleteFuncEmployeeALL)
	// Delete a department
	router.DELETE("/department/:id", handler.DeleteFuncDepartment)
	// Delete all departments
	router.DELETE("/department", handler.DeleteFuncDepartmensALL)
	router.Run(":8084")
}
