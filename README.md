A simple HTTP server that implements the storage of employees of the form 
```
type Employee struct {
	Id      int    `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	Salary  int    `json:"salary" bson:"salary"`
	Age     int    `json:"age" bson:"age"`
}
```
Implemented the ability to combine employees into departments, work with them, get statistics both for all workers in general and for individual departments.
All data is stored in the mongodb database `server` by uri `mongodb://localhost:27017` in two collections - `employees` and `department`.
Communication implemented in Rest API style by port: `localhost:8084`.
