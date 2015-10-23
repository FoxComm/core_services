package main

import (
	"log"

	"github.com/FoxComm/libs/announcer"
	"github.com/FoxComm/libs/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/core_services/router/plugins/health_check"
	"github.com/FoxComm/core_services/user/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	session := r.Group(configs.Get("UserAPIPrefix"))
	session.PUT("/login", controllers.Login)
	session.GET("/logout", controllers.Logout)

	runPort := configs.GetSafeRunPortStringFromString(configs.Get("UserPort"))

	announcer.Setup(endpoints.UserAPI.Name, runPort)
	defer announcer.Cleanup()

	health_check.Register(r)

	log.Printf("FC Router is mounting on port: %s", runPort)
	r.Run(":" + runPort)
}
