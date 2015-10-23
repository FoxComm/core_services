package controllers

import (
	"strconv"
	"time"

	"github.com/FoxComm/libs/etcd_client"
	"github.com/FoxComm/core_services/feature_manager/core"
	. "github.com/FoxComm/libs/db/masterdb"
	"github.com/gin-gonic/gin"
)

func init() {
	etcd_client.EtcdClient.Create("feature_manager_updated_at", time.Now().String(), 0)
}

func Features(c *gin.Context) {
	storeId := c.Request.Header.Get("FC-Store-ID")
	var store core.Store
	Db().First(&store, storeId)

	var features []core.Feature
	Db().Select("features.id, name, description, store_features.enabled").
		Joins("left join store_features on store_features.feature_id = features.id and store_features.store_id = " + storeId).
		Find(&features)

	c.JSON(200, features)
}

func UpdateFeatures(c *gin.Context) {
	storeIdStr := c.Request.Header.Get("FC-Store-ID")
	storeId, _ := strconv.Atoi(storeIdStr)

	var feature core.Feature
	c.Bind(&feature)

	if err := Db().Where(&core.StoreFeature{FeatureId: feature.Id, StoreId: storeId}).
		Assign("enabled", feature.Enabled).FirstOrCreate(&core.StoreFeature{}).Error; err == nil {
		etcd_client.EtcdClient.Update("feature_manager_updated_at", time.Now().String(), 0)
		feature.Name = "sss"
		c.JSON(200, feature)
	} else {
		c.JSON(500, gin.H{"error": "An error ocurred saving the record:" + err.Error()})
	}
}
