package controller

import (
	"api-hotel-booking/internal/grpc/client"
	"api-hotel-booking/internal/grpc/proto"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/responses"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStaffs(c *gin.Context) {
	q := c.Request.URL.Query()
	limitS := q.Get("limit")
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println(err)
	}

	pageS := q.Get("page")
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println(err)
	}

	id := q.Get("id")

	position := q.Get("position")

	staffs := client.GetStaffs(&proto.GetStaffsRequest{Page: int32(page), Offset: int32(limit), Id: id, Position: position}, client.GrpcClient.StaffClient)
	if staffs == nil {
		c.JSON(http.StatusOK, &proto.Staff{})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

func CreateStaff(c *gin.Context) {
	staff := models.Staff{}
	err := c.ShouldBindJSON(&staff)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	staffRequest := &proto.Staff{
		Id:          "",
		FirstName:   staff.FirstName,
		LastName:    staff.LastName,
		Position:    staff.Position,
		Salary:      staff.Salary,
		DateOfBirth: staff.DateOfBirth,
		Phone:       staff.DateOfBirth,
		Email:       staff.Email,
		StartDate:   staff.StartDate.String(),
		DeletedAt:   staff.DeletedAt.String(),
	}

	createdStaff := client.CreateStaff(staffRequest, client.GrpcClient.StaffClient)
	c.JSON(http.StatusCreated, createdStaff)
}

func UpdateStaff(c *gin.Context) {
	staff := models.Staff{}
	err := c.ShouldBindJSON(&staff)
	staff.ID = c.Param("id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	staffRequest := &proto.Staff{
		Id:          "",
		FirstName:   staff.FirstName,
		LastName:    staff.LastName,
		Position:    staff.Position,
		Salary:      staff.Salary,
		DateOfBirth: staff.DateOfBirth,
		Phone:       staff.DateOfBirth,
		Email:       staff.Email,
		StartDate:   staff.StartDate.String(),
		DeletedAt:   staff.DeletedAt.String(),
	}

	updatedStaff := client.UpdateStaff(staffRequest, client.GrpcClient.StaffClient)
	c.JSON(http.StatusOK, updatedStaff)
}

func DeleteStaff(c *gin.Context) {
	id := c.Param("id")
	staffId := &proto.DeleteStaffRequest{Id: id}

	deletedStaff := client.DeleteStaff(staffId, client.GrpcClient.StaffClient)
	c.JSON(http.StatusOK, deletedStaff)
}
