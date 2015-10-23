package main

import (
	"log"

	"github.com/FoxComm/libs/announcer"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/core_services/feature_manager/controllers"
	"github.com/FoxComm/core_services/router/plugins/health_check"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	core := r.Group(configs.Get("CoreAPIPrefix"))
	core.GET("/features", controllers.Features)
	core.PUT("/features", controllers.UpdateFeatures)

	core.GET("/stores", controllers.Store)
	core.PUT("/stores", controllers.UpdateStores)

	runPort := configs.GetSafeRunPortStringFromString(configs.Get("CorePort"))

	announcer.Setup(endpoints.CoreAPI.Name, runPort)
	defer announcer.Cleanup()

	health_check.Register(r)

	log.Printf("FC Router is mounting on port: %s", runPort)
	r.Run(":" + runPort)
}
