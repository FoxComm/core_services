package controllers

import (
	"strconv"
	"time"

	"github.com/FoxComm/FoxComm/etcd_client"
	"github.com/FoxComm/core_services/feature_manager/core"
	. "github.com/FoxComm/libs/db/masterdb"
	"github.com/gin-gonic/gin"
)

func init() {
	etcd_client.EtcdClient.Create("feature_manager_updated_at", time.Now().String(), 0)
}

func Store(c *gin.Context) {
	storeId := c.Request.Header.Get("FC-Store-ID")
	var store core.Store
	Db().First(&store, storeId)

	c.JSON(200, store)
}

func UpdateStores(c *gin.Context) {
	storeIdStr := c.Request.Header.Get("FC-Store-ID")
	storeId, _ := strconv.Atoi(storeIdStr)

	var store core.Store

	c.Bind(&store)

	if err := Db().Where(&core.Store{Id: storeId}).
		Assign("spree_token", store.SpreeToken).Assign("solr_host", store.SolrHost).FirstOrCreate(&core.Store{}).Error; err == nil {
		etcd_client.EtcdClient.Update("feature_manager_updated_at", time.Now().String(), 0)
		c.JSON(200, store)
	} else {
		c.JSON(500, gin.H{"error": "An error ocurred saving the record:" + err.Error()})
	}
}
