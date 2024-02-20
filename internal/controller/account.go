package controller

import (
	"api-hotel-booking/internal/database"
	"api-hotel-booking/internal/grpc/client"
	"api-hotel-booking/internal/grpc/proto"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/repository"
	"api-hotel-booking/internal/responses"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
	q := c.Request.URL.Query()
	limitS := q.Get("limit")
	pageS := q.Get("page")
	//log.Fatal(q)
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println(err)

	}
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println(err)

	}
	accounts := client.GetAllAccount(proto.AccountRequest{Page: int32(page), Offset: int32(limit)}, client.GrpcClient.AccountClient)
	c.JSON(http.StatusOK, accounts)
}

func CreateAccount(c *gin.Context) {
	account := models.Account{}
	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	account.Prepare()
	err = account.Validtate("create")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := repository.NewRepositoryAccountCRUD(db)
	func(accountRepository repository.AccountRepo) {
		account, err = accountRepository.Save(account)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		//w.Header().Set("Location", fmt.Sprintf("%s %s %d", r.Host, r.RequestURI, account.ID))
		c.JSON(http.StatusCreated, account)
	}(repo)
}
func UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.ERROR(http.StatusBadRequest, err.Error())
	}

	account := models.Account{}
	err = c.ShouldBindJSON(&account)
	fmt.Println(account)
	account.Prepare()
	err = account.Validtate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := repository.NewRepositoryAccountCRUD(db)

	func(accountRepository repository.AccountRepo) {
		rows, err := accountRepository.Update(uint32(uid), account)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d row affected", rows))
	}(repo)
}
func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.ERROR(http.StatusBadRequest, err.Error())
	}
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := repository.NewRepositoryAccountCRUD(db)

	func(accountRepository repository.AccountRepo) {
		rows, err := accountRepository.Delete(uint32(uid))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d row deleted", rows))
	}(repo)
}
func FindByEmail(c *gin.Context) {
	q := c.Request.URL.Query()
	email := q.Get("email")
	fmt.Println("email:" + email)
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := repository.NewRepositoryAccountCRUD(db)
	//pagination := []models.Account{}
	func(accountRepository repository.AccountRepo) {
		account, err := accountRepository.FindByEmail(email)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, account)
	}(repo)
}
func FindByMerchantCode(c *gin.Context) {
	q := c.Request.URL.Query()
	merchantcode := q.Get("merchantcode")
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := repository.NewRepositoryAccountCRUD(db)
	//pagination := []models.Account{}
	func(accountRepository repository.AccountRepo) {
		accounts, err := accountRepository.FindByMerchantCode(merchantcode)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, accounts)
	}(repo)
}
