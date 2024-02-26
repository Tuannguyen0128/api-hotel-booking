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

func GetAccounts(c *gin.Context) {
	q := c.Request.URL.Query()
	limitS := q.Get("limit")
	pageS := q.Get("page")
	id := q.Get("id")
	username := q.Get("username")
	staffID := q.Get("staff_id")
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println(err)
	}
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println(err)
	}
	accounts := client.GetAccounts(&proto.GetAccountsRequest{Page: int32(page), Offset: int32(limit), Id: id, Username: username, StaffId: staffID}, client.GrpcClient.AccountClient)
	if accounts == nil {
		c.JSON(http.StatusOK, &proto.Account{})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func CreateAccount(c *gin.Context) {
	account := models.Account{}
	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	accountRequest := &proto.Account{
		Id:         "",
		StaffId:    account.StaffID,
		Username:   account.Username,
		Password:   account.Password,
		UserRoleId: account.UserRoleID,
	}

	createdAccount := client.CreateAccount(accountRequest, client.GrpcClient.AccountClient)
	c.JSON(http.StatusCreated, createdAccount)
}

func UpdateAccount(c *gin.Context) {
	account := models.Account{}
	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	accountRequest := &proto.Account{
		Id:         c.Param("id"),
		StaffId:    account.StaffID,
		Username:   account.Username,
		Password:   account.Password,
		UserRoleId: account.UserRoleID,
	}

	updatedAccount := client.UpdateAccount(accountRequest, client.GrpcClient.AccountClient)
	c.JSON(http.StatusOK, updatedAccount)
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	accountId := &proto.DeleteAccountRequest{Id: id}

	deletedAccount := client.DeleteAccount(accountId, client.GrpcClient.AccountClient)
	c.JSON(http.StatusOK, deletedAccount)
}
