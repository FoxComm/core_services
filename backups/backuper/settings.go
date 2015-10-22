package backuper

import (
	"github.com/FoxComm/FoxComm/utils"
	. "github.com/FoxComm/libs/db/masterdb"
	"github.com/gin-gonic/gin"
)

type Settings struct {
	Id            int
	StoreId       int
	S3Id          string
	S3Secret      string
	S3Bucket      string
	S3Region      string
	DropboxId     string
	DropboxSecret string
	DropboxToken  string
	DriveId       string
	DriveSecret   string
	DatabaseUrl   string
	AutoBackup    bool
}

func MountSettings(r *gin.RouterGroup) {
	settings := r.Group("settings", SettingsBase)
	settings.GET("/", SettingsShow)
	settings.PUT("/:id", SettingsUpdate)
}

func (s Settings) TableName() string {
	return "backup_settings"
}

func SettingsBase(c *gin.Context) {
	c.Next()
}

func SettingsShow(c *gin.Context) {
	storeId := utils.StoreID(c)
	settings := Settings{}
	Db().FirstOrCreate(&settings, Settings{StoreId: storeId})
	c.JSON(200, &settings)
}

func SettingsUpdate(c *gin.Context) {
	storeId := utils.StoreID(c)
	settings := Settings{}
	if c.Bind(&settings) && settings.StoreId == storeId {
		Db().Save(&settings)
		c.JSON(200, &settings)
	} else {
		c.JSON(422, gin.H{"Error": "Can't process request"})
	}
}
