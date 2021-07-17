package server

import (
	"io"
	"log"
	"os"

	_ "github.com/Monkey-Mouse/mo2/docs"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/mo2utils/basiclog"
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

	if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()
		f, _ := os.OpenFile("logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		errf, _ := os.OpenFile("logs/err.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		gin.DefaultWriter = io.MultiWriter(f)
		gin.DefaultErrorWriter = io.MultiWriter(errf)
		basiclog.SetLoggeer(log.New(f, "[INFO]", 0), log.New(errf, "[ERROR]", 0))
	}
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	c := controller.NewController()
	controller.SetupHandlers(c)
	fromCTX := func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
		str, err := ctx.Cookie("jwtToken")
		if err != nil {
			return
		}
		userInfo, err = mo2utils.ParseJwt(str)
		return
	}
	middleware.H.RegisterMapedHandlers(r,
		&middleware.OptionalParams{
			LimitEvery:     10,
			Unblockevery:   3600,
			UseRedis:       true,
			GetUserFromCTX: fromCTX,
			UserKey:        mo2utils.UserInfoKey,
		})
	pprof.Register(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "dist/index.html")
	})
	r.Run(":5001")
}
