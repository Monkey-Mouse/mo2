package mo2utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserName string `json:"user_name"`
	Login    bool   `json:"login"`
	Infos    []byte `json:"infos"`
}

type JwtClaims struct {
	UserInfo UserInfo `json:"user_info"`
	jwt.StandardClaims
}

func (j JwtClaims) Valid() error {

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
		if userInfo.Login {
			c.JSON(http.StatusOK, gin.H{"user": userInfo.UserName})
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
func VerifyJwt(tokenString string) (userInfo UserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mo2"), nil
	})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		userInfo = claims.UserInfo
		err = nil
	} else {
		log.Println("can not parse with JwtClaims")
	}
	return

}

//add jwt to generate token for user
//if token is valid, return nil
func ParseJwt(tokenString string) (userInfo UserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mo2"), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		userInfo = claims.UserInfo
	} else {
		log.Println("can not parse with JwtClaims")
	}
	return

}

//add jwt to generate token for user
//if token is valid, return nil
func verifyJwt2(tokenString string, claims interface{}) (err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//TODO to change the key
		return []byte("shhhh"), nil
	})
	if token.Valid {
		claims := token.Claims

		fmt.Println(claims)
		//jwt.DecodeSegment()
		// process to decode by base64 the info in jwt token
		fmt.Println("Valid!")
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Token is either expired or not active yet")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	return
}

// generate jwt token with claim in type JwtClaims
// for info
// jwtToken string
func GenerateJwtCode(info string, infos interface{}) string {
	//TODO change the key
	hmacSampleSecret := []byte("mo2")
	infosByte, err := json.Marshal(infos)
	if err != nil {
		log.Fatal(err)
	}
	claims := JwtClaims{
		UserInfo: UserInfo{
			UserName: info,
			Login:    true,
			Infos:    infosByte,
		},
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

func VerifyJwtExample() {
	//	var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
	var tokenString = "JhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.gB7MbL6QovrUe1d7AJIo_MZ5NZmIp30g7eeG8ZeQWz8"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("shhhh"), nil
	})
	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
}
