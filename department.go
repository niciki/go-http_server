package main

import (
	"fmt"
)

type Department struct {
	CountStaff   int        `json:"countstaff"`
	Staff        []Employee `json:"staff"`
	DepartmentID int        `json:"department_id"`
}

type DepartmentID struct {
	CountStaff   int   `json:"countstaff" bson:"countstaff"`
	StaffID      []int `json:"staffid" bson:"staffid"`
	DepartmentID int   `json:"department_id" bson:"department_id"`
}

func FillDepartment(depid *DepartmentID, db *Database) (Department, error) {
	var dep Department
	dep.CountStaff = depid.CountStaff
	dep.DepartmentID = depid.DepartmentID
	staff := make([]Employee, len(depid.StaffID))
	db.m.RLock()
	defer db.m.RUnlock()
	for i, j := range depid.StaffID {
		data, err := db.GetID(j)
		if err != nil {
			return dep, fmt.Errorf("there isn't an id: {%d} in database", i)
		} else {
			staff[i] = data
		}
	}
	dep.Staff = staff
	return dep, nil
}
