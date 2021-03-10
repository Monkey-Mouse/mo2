package main

import (
	"mo2/docs"
	"mo2/mo2utils"
	"mo2/server"
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

// @host localhost:5001
// @BasePath /
// @query.collection.format multi
func main() {
	server.RunServer()
}

func init() {
	if mo2utils.IsEnvRelease() {
		docs.SwaggerInfo.Version = "v0.2 alpha"
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = "www.motwo.cn"
	}
	mo2utils.UploadCDN()
}
