package controller

import (
	"ProjectPractice/src/api/database"
	"ProjectPractice/src/api/models"
	"ProjectPractice/src/api/repository"
	"ProjectPractice/src/api/repository/crud"
	"ProjectPractice/src/api/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetMerchantAccounts(c *gin.Context) {
	var limit, page int
	var err error
	q := c.Request.URL.Query()
	limitS := q.Get("limit")
	pageS := q.Get("page")
	//log.Fatal(q)
	if limitS != "" {
		limit, err = strconv.Atoi(limitS)

	}
	if pageS != "" {
		limit, err = strconv.Atoi(pageS)

	}
	if err != nil {
		log.Println(err)

	}
	page, err = strconv.Atoi(pageS)
	if err != nil {
		log.Println(err)

	}
	pagination := models.Pagination{Limit: limit, Page: page}

	//err = json.Unmarshal(body, &pagination)
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryMerchantAccountCRUD(db)
	func(merchantAccountRepository repository.MerchantAccountRepository) {
		pagination, err = merchantAccountRepository.FindAll(pagination)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, pagination)
	}(repo)
}
func CreateMerchantAccount(c *gin.Context) {

	merchantaccount := models.MerchantAccount{}
	err := c.ShouldBindJSON(&merchantaccount)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	fmt.Println(merchantaccount)
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryMerchantAccountCRUD(db)
	func(merchantAccountRepository repository.MerchantAccountRepository) {
		merchantaccount, err = merchantAccountRepository.Save(merchantaccount)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		//c.Writer.Header().Set("Location", fmt.Sprintf("%s %s %d", r.Host, r.RequestURI, merchantaccount.ID))
		c.JSON(http.StatusCreated, merchantaccount)
	}(repo)
}
func UpdateMerchantAccount(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,responses.ERROR(http.StatusBadRequest, err.Error()))
	}
	merchantaccount := models.MerchantAccount{}
	err = c.ShouldBindJSON(&merchantaccount)
	fmt.Println(merchantaccount)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryMerchantAccountCRUD(db)

	func(merchantAccountRepository repository.MerchantAccountRepository) {
		rows, err := merchantAccountRepository.Update(uint32(uid), merchantaccount)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d row affected", rows))
	}(repo)
}
func FindMerchantAccountByID(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,responses.ERROR(http.StatusBadRequest, err.Error()))
	}
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryMerchantAccountCRUD(db)
	func(merchantAccountRepository repository.MerchantAccountRepository) {
		teammembers, err := merchantAccountRepository.FindByID(uint32(uid))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON( http.StatusOK, teammembers)
	}(repo)
}
func DeleteMerchantAccount(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,responses.ERROR(http.StatusBadRequest, err.Error()))
	}
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryMerchantAccountCRUD(db)

	func(merchantAccountRepository repository.MerchantAccountRepository) {
		rows, err := merchantAccountRepository.Delete(uint32(uid))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d row affected", rows))
	}(repo)
}
