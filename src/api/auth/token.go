package auth

import (
	"ProjectPractice/src/api/utils/console"
	"ProjectPractice/src/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(user_id string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.SECRETKEY)
}
func TokenValidate(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.SECRETKEY, nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		console.Pretty(claims)
		return nil
	} else {
		fmt.Println(err)
		return err
	}

}
func ExtractToken(c *gin.Context) string {

	bearerToken := c.Request.Header.Get("Authorization")
	//println("Bearer token:" + bearerToken)
	log.Println("Bearer token:" + bearerToken)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
