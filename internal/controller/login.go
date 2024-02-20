package controller

import (
	"api-hotel-booking/internal/auth"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	teammember := models.Account{}
	err := c.ShouldBindJSON(&teammember)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	teammember.Prepare()
	err = teammember.Validtate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	tokenrp, err := auth.SignIn(teammember.Email, teammember.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	c.JSON(http.StatusAccepted, models.TokenResponse{Token: tokenrp, Expire_In: time.Now().Add(time.Minute * 15).Unix(), Token_Type: "bearer"})
}
