package main

import (
	_ "mo2/docs"
	"mo2/server"
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
	server.RunServer()

	//demo.GetClient()

	//demo.Welcome()

}
