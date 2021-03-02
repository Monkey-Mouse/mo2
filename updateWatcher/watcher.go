package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/post", func(ctx *gin.Context) {
		dir, _ := os.Getwd()
		fmt.Println("enter", dir)
		ex := exec.Command("/bin/bash", "./refresh.bash")
		fmt.Println("exec")
		cmd, err := ex.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(cmd))
		ctx.JSON(http.StatusOK, nil)
	})
	r.Run(":5002")
}
