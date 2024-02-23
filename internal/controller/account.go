package controller

import (
	"api-hotel-booking/internal/grpc/client"
	"api-hotel-booking/internal/grpc/proto"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/responses"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
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
	accounts := client.GetAllAccount(&proto.GetAccountsRequest{Page: int32(page), Offset: int32(limit), Id: id, Username: username, StaffId: staffID}, client.GrpcClient.AccountClient)
	c.JSON(http.StatusOK, accounts)
}

func CreateAccount(c *gin.Context) {
	account := models.Account{}
	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
	}

	accountRequest := &proto.Account{
		Id:          id.String(),
		StaffId:     account.StaffID,
		Username:    account.Username,
		Password:    account.Password,
		UserRoleId:  account.UserRoleID,
		CreatedAt:   account.CreatedAt.String(),
		UpdatedAt:   account.UpdatedAt.String(),
		DeletedAt:   account.DeletedAt.String(),
		LastLoginAt: account.LastLoginAt.String(),
	}

	createdAccount := client.CreateAccount(accountRequest, client.GrpcClient.AccountClient)
	c.JSON(http.StatusCreated, createdAccount)
}

func UpdateAccount(c *gin.Context) {
	account := models.Account{}
	err := c.ShouldBindJSON(&account)
	account.ID = c.Param("id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	accountRequest := &proto.Account{
		Id:          account.ID,
		StaffId:     account.StaffID,
		Username:    account.Username,
		Password:    account.Password,
		UserRoleId:  account.UserRoleID,
		CreatedAt:   account.CreatedAt.String(),
		UpdatedAt:   account.UpdatedAt.String(),
		DeletedAt:   account.DeletedAt.String(),
		LastLoginAt: account.LastLoginAt.String(),
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
