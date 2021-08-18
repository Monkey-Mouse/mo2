package mo2utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils/redisutil"

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

const keystr = "MO2_KEY"

// GetKey from redis or create new key
func GetKey() (key []byte, err error) {
	key = make([]byte, 16)
	red := redisutil.GetRedisClient()
	skey, err := red.Get(keystr).Result()
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		_, err = rand.Read(key)
		if err != nil {
			return nil, err
		}
		suc, err := red.SetNX(keystr, string(key), time.Hour*2400).Result()
		if err != nil {
			return nil, err
		}
		if !suc {
			skey, err = red.Get(keystr).Result()
			if err != nil {
				return nil, err
			}
			key = []byte(skey)
		}
		return key, nil
	}
	key = []byte(skey)
	return
}

//ParseJwt for JwtLoginClaims
//if token is valid, return nil
func ParseJwt(tokenString string) (userInfo dto.LoginUserInfo, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtLoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetKey()
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
	key, err := GetKey()
	if err != nil {
		panic(err)
	}
	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString
}

// GenerateVerifyJwtToken generate jwt token with claim in type JwtClaims
// for info string
// jwtToken string
func GenerateVerifyJwtToken(info string) string {
	return GenerateJwtToken(info, time.Now().Add(time.Minute*10))
}

// GenerateJwtToken generate jwt token with claim in type JwtClaims
// for info string
// jwtToken string
func GenerateJwtToken(info string, expireAt time.Time) string {

	claims := JwtInfoClaims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	key, err := GetKey()
	if err != nil {
		panic(err)
	}
	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString
}
