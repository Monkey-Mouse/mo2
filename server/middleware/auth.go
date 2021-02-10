package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BasicAuth() gin.HandlerFunc {

	//check the cookie
	return func(c *gin.Context) {

		login, err := c.Cookie("login")
		if login == "true" {
			user, err := c.Cookie("user")
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, gin.H{"login": login, "user": user})
		} else {
			log.Println(err)
			//TODO redirect to login
			//display login window to get post
			c.JSON(http.StatusForbidden, gin.H{"info": "need to login first"})
		}
	}
}
