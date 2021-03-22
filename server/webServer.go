package server

import (
	_ "github.com/Monkey-Mouse/mo2/docs"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/controller"
	"github.com/Monkey-Mouse/mo2/server/middleware"

	"github.com/gin-contrib/pprof"

	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RunServer start web server
func RunServer() {

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	c := controller.NewController()
	controller.SetupHandlers(c)
	middleware.H.RegisterMapedHandlers(r, func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
		str, err := ctx.Cookie("jwtToken")
		if err != nil {
			return
		}
		userInfo, err = mo2utils.ParseJwt(str)
		return
	}, mo2utils.UserInfoKey, &middleware.OptionalParams{LimitEvery: 10, Unblockevery: 3600, UseRedis: true})
	pprof.Register(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "dist/index.html")
	})
	r.Run(":5001")
}
