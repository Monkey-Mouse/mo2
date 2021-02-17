package main

import (
	"mo2/config"
	_ "mo2/docs"
	demo "mo2/examples"
	"os"
	//"time"
)

// @title Mo2
// @version 1.0
// @description This is a Mo2 server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
// @query.collection.format multi
func main() {
	//expTime:=time.Date(2021,2,14,12,0,0,0,time.UTC).Unix()

	//middleware.GenerateJwtCode()
	demo.SendEmail()
	//server.RunServer()

	//demo.GetClient()

	//demo.Welcome()

}
func init() {
	//config.SaveConfig("../config/config.yaml")
	c := config.LoadConfig("../config/config.yaml")
	os.Setenv("emailAddr", c.EmailAddr)
	os.Setenv("emailPass", c.EmailPass)
}
