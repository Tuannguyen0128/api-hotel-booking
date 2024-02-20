package controller

import (
	"ProjectPractice/src/api/auth"
	"ProjectPractice/src/api/models"
	"ProjectPractice/src/api/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	teammember := models.TeamMember{}
	err := c.ShouldBindJSON(&teammember)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	teammember.Prepare()
	err = teammember.Validtate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	tokenrp, err := auth.SignIn(teammember.Email, teammember.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	c.JSON(http.StatusAccepted, models.TokenResponse{Token: tokenrp, Expire_In: time.Now().Add(time.Minute * 15).Unix(), Token_Type: "bearer"})
}
