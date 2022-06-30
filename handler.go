package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db      Database
	company Company
}

func NewHandler(collectiondb *mongo.Collection, collectioncom *mongo.Collection) *Handler {
	return &Handler{db: *NewDatabase(collectiondb), company: Company{length: 0, collection: collectioncom}}
}

func (h *Handler) PostFuncEmployee(c *gin.Context) {
	var empl Employee
	if err := c.BindJSON(&empl); err != nil {
		fmt.Printf("Fail during bind employee from json: %s\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.db.Add(&empl)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"size": h.db.length,
	})
}

func (h *Handler) GetFuncEmployeeID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Fail during converting id to int: %s\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		val, err := h.db.GetID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"id": val,
			})
		}
	}
}

func (h *Handler) GetFuncEmployeeALL(c *gin.Context) {
	data, err := h.db.GetAll()
	if err == nil {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) PutFuncEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Fail during converting id to int: %s\n", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"id": err.Error(),
		})
	} else {
		var empl Employee
		if err := c.BindJSON(&empl); err != nil {
			fmt.Printf("Fail during bind employee from json: %s\n", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"id": err.Error(),
			})
		}
		h.db.Put(id, &empl)
		c.JSON(http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

func (h *Handler) DeleteFuncEmployeeID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Fail during converting id to int: %s\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
	} else if _, err := h.db.GetID(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		h.db.DeleteID(id)
		c.JSON(http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

func (h *Handler) DeleteFuncEmployeeALL(c *gin.Context) {
	num := h.db.length
	err := h.db.DeleteALL()
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"number of delete elements": num,
		})
	} else {
		c.JSON(http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) PostFuncDepartment(c *gin.Context) {
	var dep DepartmentID
	if err := c.BindJSON(&dep); err != nil {
		fmt.Printf("Fail during bind DepartmentID from json: %s\n", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"id": err.Error(),
		})
	}
	err := h.company.Add(&dep)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"company_size": h.company.length,
		})
	}

}

func (h *Handler) GetFuncDepartmentID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("fail during converting id to int: %s\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		val, err := h.company.GetID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"id": val,
			})
		}
	}
}

func (h *Handler) GetFuncDepartmentALL(c *gin.Context) {
	data, err := h.company.GetALL(&h.db)
	if err == nil {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) PutFuncAddDepartment(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Param(("id")))
	if err1 != nil {
		fmt.Printf("fail during converting id to int: %s\n", err1)
		c.JSON(http.StatusBadRequest, err1.Error())
	} else {
		_, err := h.company.GetID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	employeeid, err3 := strconv.Atoi(c.Param(("employeeid")))
	if err3 != nil {
		fmt.Printf("fail during converting employeeid to int: %s\n", err3)
		c.JSON(http.StatusBadRequest, err3.Error())
	} else {
		_, err4 := h.db.GetID(employeeid)
		if err4 != nil {
			c.JSON(http.StatusBadRequest, "there isn't employee with such id")
			return
		}
	}
	dep, _ := h.company.GetID(id)
	for _, j := range dep.StaffID {
		if j == employeeid {
			c.JSON(http.StatusOK, "this employee already exists")
			return
		}
	}
	dep.StaffID = append(dep.StaffID, employeeid)
	dep.CountStaff++
	h.company.Put(id, &dep)
	c.JSON(http.StatusOK, dep)
}

func (h *Handler) PutFuncDeleteDepartment(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Param(("id")))
	if err1 != nil {
		fmt.Printf("fail during converting id to int: %s\n", err1)
		c.JSON(http.StatusBadRequest, err1.Error())
	} else {
		_, err := h.company.GetID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	employeeid, err3 := strconv.Atoi(c.Param(("employeeid")))
	if err3 != nil {
		fmt.Printf("fail during converting employeeid to int: %s\n", err3)
		c.JSON(http.StatusBadRequest, err3.Error())
	} else {
		_, err4 := h.db.GetID(employeeid)
		if err4 != nil {
			c.JSON(http.StatusBadRequest, "there isn't employee with such id")
			return
		}
	}
	dep, _ := h.company.GetID(id)
	staffid := make([]int, 0)
	for _, j := range dep.StaffID {
		if j != employeeid {
			staffid = append(staffid, j)
		}
	}
	dep.StaffID = staffid
	h.company.Put(id, &dep)
	c.JSON(http.StatusOK, dep)
}

func (h *Handler) DeleteFuncDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Fail during converting id to int: %s\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
	} else if _, err := h.company.GetID(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		h.company.DeleteID(id)
		c.JSON(http.StatusOK, id)
	}
}

func (h *Handler) GetStat(c *gin.Context) {
	val, err := NewStatistic(h)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, val)
}

func (h *Handler) DeleteFuncDepartmensALL(c *gin.Context) {
	err := h.company.DeleteALL()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, "Delete all departments successfully")
}
