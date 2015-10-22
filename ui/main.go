package main

import (
	"net/http"
	"net/url"

	"github.com/FoxComm/FoxComm/announcer"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/libs/logger"
)

func main() {
	logger.Debug("Starting UI Server")
	if _, err := url.Parse(configs.Get("UIHost")); err == nil {
		//		http.Handle("/images", http.FileServer(http.Dir("/dist/images")))
		//		http.Handle("/styles", http.FileServer(http.Dir("/dist/styles")))

		http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			logger.Debug("Serving index.html")
			http.ServeFile(response, request, "/dist/index.html")
		})

		port := configs.GetSafeRunPortStringFromString(configs.Get("UIPort"))

		announcer.Setup(endpoints.UIAPI.Name, port)
		defer announcer.Cleanup()
		logger.Debug("Starting server in port:" + port)
		http.ListenAndServe(":"+port, http.FileServer(http.Dir("dist")))

	} else {
		logger.Error("UIHost is not set")
	}
}
