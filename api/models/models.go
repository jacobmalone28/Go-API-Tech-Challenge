package models

type Course struct {
	ID int
	Name string
}

type Person struct {
	ID 			int
	FirstName 	string
	LastName 	string
	Type 		string
	Age 		int
	Courses 	[]int
}