package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jrallison/go-workers"

	"github.com/FoxComm/FoxComm/announcer"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/core_services/backups/backuper"
	"github.com/FoxComm/core_services/router/plugins/health_check"
)

func main() {
	r := gin.Default()

	backups := r.Group(endpoints.BackupsAPI.APIPrefix, BackupsBase)

	backuper.MountSettings(backups)
	backuper.MountJobs(backups)

	backups.GET("stats", stats)

	backuper.Configure()

	runPort := configs.GetSafeRunPortStringFromString(endpoints.BackupsAPI.DefaultPort)
	logger.Info("Backup service is mounting on port: %s", runPort)

	health_check.Register(r)

	announcer.Setup(endpoints.BackupsAPI.Name, runPort)
	defer announcer.Cleanup()

	go r.Run("0.0.0.0:" + runPort)

	backuper.Start(10)
}

func BackupsBase(c *gin.Context) {
	c.Next()
}

func stats(c *gin.Context) {
	workers.Stats(c.Writer, c.Request)
}
