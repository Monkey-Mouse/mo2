package mo2utils

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
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
type JwtInfoClaims struct {
	Info string `json:"info"`
	jwt.StandardClaims
}

var key []byte = make([]byte, 16)

func init() {
	bytes, err := ioutil.ReadFile("mo2.secret")
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		rand.Read(key)
		ioutil.WriteFile("mo2.secret", key, 0)
		return
	}
	key = bytes
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

// VerifyInfoJwt for JwtInfoClaims
//if token is valid, return nil
func VerifyInfoJwt(tokenString string) (info string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mo2"), nil
	})

	if claims, ok := token.Claims.(*JwtInfoClaims); ok && token.Valid {
		info = claims.Info
		err = nil
	} else {
		err = errors.New("can not parse for jwtInfoClaims")
	}
	return

}

//ParseJwt for JwtLoginClaims
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
		err = errors.New("can not parse for jwtLoginClaims")
	}
	return

}

// GenerateJwtCode generate jwt token with claim in type JwtClaims
// for info dto.LoginUserInfo
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
	return tokenString
}

// GenerateJwtToken generate jwt token with claim in type JwtClaims
// for info string
// jwtToken string
func GenerateJwtToken(info string) string {
	//TODO change the key

	claims := JwtInfoClaims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
