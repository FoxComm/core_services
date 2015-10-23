package main

import (
	"github.com/FoxComm/libs/announcer"
	"github.com/FoxComm/core_services/catalog_cache/controllers"
	"github.com/FoxComm/libs/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/core_services/router/plugins/health_check"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	catalog := r.Group(configs.Get("CatalogCacheAPIPrefix"))
	catalog.GET("/products", controllers.Products)
	catalog.GET("/products/:slug", controllers.Product)

	runPort := configs.GetSafeRunPortStringFromString(configs.Get("CatalogCachePort"))
	logger.Debug("CatalogCache is mounting on port: %s", runPort)

	health_check.Register(r)
	announcer.Setup(endpoints.CatalogCacheAPI.Name, runPort)
	defer announcer.Cleanup()

	r.Run(":" + runPort)
}
