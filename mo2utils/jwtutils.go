package mo2utils

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"mo2/dto"
	"os"
	"path"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtLoginClaims mo2 jwt claims
type JwtLoginClaims struct {
	UserInfo dto.LoginUserInfo `json:"user_info"`
	jwt.StandardClaims
}

// JwtInfoClaims claim with extra string info
type JwtInfoClaims struct {
	Info string `json:"info"`
	jwt.StandardClaims
}

var key []byte = make([]byte, 16)

func init() {
	initKey()
}
func initKey() (err error) {
	os.Mkdir("./secrets", 0755)
	path := path.Join("./secrets", "mo2.secret")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {

		rand.Seed(time.Now().UnixNano())
		_, err = rand.Read(key)
		err = ioutil.WriteFile(path, key, 0)
		return
	}
	key = bytes
	return
}

//ParseJwt for JwtLoginClaims
//if token is valid, return nil
func ParseJwt(tokenString string) (userInfo dto.LoginUserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtLoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
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

	claims := JwtLoginClaims{
		UserInfo: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
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

// GenerateJwtToken generate jwt token with claim in type JwtClaims
// for info string
// jwtToken string
func GenerateJwtToken(info string) string {

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
