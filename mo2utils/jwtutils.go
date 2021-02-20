package mo2utils

import (
	"fmt"
	"log"
	"mo2/dto"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtLoginClaims struct {
	UserInfo dto.LoginUserInfo `json:"user_info"`
	jwt.StandardClaims
}

func (j JwtLoginClaims) Valid() error {

	return nil

}

func BasicAuth() gin.HandlerFunc {

	//check the cookie
	return func(c *gin.Context) {

		jwtToken, err := c.Cookie("jwtToken")
		if err != nil {
			log.Println(err)
			//c.JSON(http.StatusForbidden, gin.H{"info": "need to login first"})

		}
		userInfo, err := VerifyJwt(jwtToken)
		c.Keys[""] = ""
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"user": userInfo.Name})
		} else {
			log.Println(err)
			//TODO redirect to login
			//display login window to get post
			c.JSON(http.StatusForbidden, gin.H{"info": "need to login first"})
		}
	}
}

//add jwt to generate token for user
//if token is valid, return nil
func VerifyJwt(tokenString string) (userInfo dto.LoginUserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtLoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mo2"), nil
	})

	if claims, ok := token.Claims.(*JwtLoginClaims); ok && token.Valid {
		userInfo = claims.UserInfo
		err = nil
	} else {
		log.Println("can not parse with JwtClaims")
	}
	return

}

//add jwt to generate token for user
//if token is valid, return nil
func ParseJwt(tokenString string) (userInfo dto.LoginUserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtLoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mo2"), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*JwtLoginClaims); ok && token.Valid {
		userInfo = claims.UserInfo
	} else {
		log.Println("can not parse with JwtClaims")
	}
	return

}

// generate jwt token with claim in type JwtClaims
// for info
// jwtToken string
func GenerateJwtCode(info dto.LoginUserInfo) string {
	//TODO change the key
	hmacSampleSecret := []byte("mo2")

	claims := JwtLoginClaims{
		UserInfo: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tokenString)
	return tokenString
}
