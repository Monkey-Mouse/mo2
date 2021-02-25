package controller

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

// Deploy will be called after newer docker img built
func (c *Controller) Deploy(ctx *gin.Context) {
	cmd := exec.Command("/bin/bash", "docker-compose --env-file ~/mo2/var.env pull&&docker-compose --env-file ~/mo2/var.env up -d")
	cmd.Run()
	cmd.Wait()
}
