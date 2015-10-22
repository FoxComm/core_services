package backuper

import (
	"github.com/gin-gonic/gin"
	"github.com/jrallison/go-workers"

	"github.com/FoxComm/FoxComm/logger"
	"github.com/FoxComm/FoxComm/utils"
	. "github.com/FoxComm/libs/db/masterdb"
)

type Job struct {
	Action   string
	Settings Settings
}

type ActionsRequest struct {
	Actions []string `json:"actions"`
}

func MountJobs(r *gin.RouterGroup) {
	jobs := r.Group("jobs")
	jobs.POST("/", CreateJob)
}

func CreateJob(c *gin.Context) {
	storeID := utils.StoreID(c)
	settings := Settings{StoreId: storeID}

	err := Db().First(&settings, settings).Error

	if err != nil {
		logger.Error("[backuper] CreateJob(): %s", err.Error())
	}

	actionsRequest := ActionsRequest{}

	if !c.Bind(&actionsRequest) {
		c.AbortWithStatus(500)
	}

	logger.Info("Backup actions: %+v", actionsRequest)

	for _, action := range actionsRequest.Actions {
		job := Job{
			Action:   action,
			Settings: settings,
		}
		workers.Enqueue("backups", "Add", job)
	}

	c.JSON(201, gin.H{"status": "enqueued"})
}
