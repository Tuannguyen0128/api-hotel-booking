package controller

import (
	"api-hotel-booking/internal/database"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/repository"
	"api-hotel-booking/internal/repository/crud"
	"api-hotel-booking/internal/responses"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTeamMembers(c *gin.Context) {

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
	pagination := models.Pagination{Limit: limit, Page: page}

	//err = json.Unmarshal(body, &pagination)
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryTeamMemberCRUD(db)
	//pagination := []models.TeamMember{}
	func(teamMemberRepository repository.TeamMemberRepository) {
		pagination, err = teamMemberRepository.FindAll(pagination)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, pagination)
	}(repo)
}
func GetTeamMember(c *gin.Context) {

}
func CreateTeamMember(c *gin.Context) {
	teammember := models.TeamMember{}
	err := c.ShouldBindJSON(&teammember)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	teammember.Prepare()
	err = teammember.Validtate("create")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	repo := crud.NewRepositoryTeamMemberCRUD(db)
	func(teamMemberRepository repository.TeamMemberRepository) {
		teammember, err = teamMemberRepository.Save(teammember)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		//w.Header().Set("Location", fmt.Sprintf("%s %s %d", r.Host, r.RequestURI, teammember.ID))
		c.JSON(http.StatusCreated, teammember)
	}(repo)
}
func UpdateTeamMember(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.ERROR(http.StatusBadRequest, err.Error())
	}

	teammember := models.TeamMember{}
	err = c.ShouldBindJSON(&teammember)
	fmt.Println(teammember)
	teammember.Prepare()
	err = teammember.Validtate("update")
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
	repo := crud.NewRepositoryTeamMemberCRUD(db)

	func(teamMemberRepository repository.TeamMemberRepository) {
		rows, err := teamMemberRepository.Update(uint32(uid), teammember)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d row affected", rows))
	}(repo)
}
func DeleteTeamMember(c *gin.Context) {
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
	repo := crud.NewRepositoryTeamMemberCRUD(db)

	func(teamMemberRepository repository.TeamMemberRepository) {
		rows, err := teamMemberRepository.Delete(uint32(uid))
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
	repo := crud.NewRepositoryTeamMemberCRUD(db)
	//pagination := []models.TeamMember{}
	func(teamMemberRepository repository.TeamMemberRepository) {
		teammember, err := teamMemberRepository.FindByEmail(email)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, teammember)
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
	repo := crud.NewRepositoryTeamMemberCRUD(db)
	//pagination := []models.TeamMember{}
	func(teamMemberRepository repository.TeamMemberRepository) {
		teammembers, err := teamMemberRepository.FindByMerchantCode(merchantcode)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		c.JSON(http.StatusOK, teammembers)
	}(repo)
}
