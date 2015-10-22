package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/FoxComm/FoxComm/announcer"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/FoxComm/router/common"
	"github.com/FoxComm/FoxComm/router/netutils"
	"github.com/FoxComm/FoxComm/router/registry"
	"github.com/FoxComm/vulcand/log"
	"github.com/FoxComm/vulcand/plugin"
	"github.com/FoxComm/vulcand/service"
)

const (
	vulcanFuncRegexp = `(?i)FoxComm\/vulcand\/(\w+.+)\.(.+)$`
	vulcanFileRegexp = `(?i)FoxComm\/vulcand\/(\w+.+\.\w+)$`
	packRegexp       = `(^\w+)\.(.+)`
)

func init() {
	if configs.Get("FC_ENV") == "development" {
		announcer.Setup(endpoints.OriginFrontendAPI.Name, "8080") //Not sure if this should live here
		announcer.Setup(endpoints.OriginBackendAPI.Name, "8080")  //Not sure if this should live here
		announcer.Setup(endpoints.UIAPI.Name, "9001")
	}
}

type NoSrvHandler struct {
}

func (e *NoSrvHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, err error) {
	statusCode := http.StatusServiceUnavailable
	netutils.WriteJsonResponse(w, statusCode, map[string]string{"error": http.StatusText(statusCode)})
}

func Run(registry *plugin.Registry) error {
	options, err := service.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("failed to parse command line: %s", err)
	}

	registry.SetNoServersErrorHandler(&NoSrvHandler{})

	if configs.Get("FC_ENV") == "development" {
		options.TomlConfigPaths.Set(configs.Get("TomlServersDir"))
		options.TomlWatchConfigChanges = true
	}

	srv := service.NewService(options, registry)
	logger, err := common.NewLogger(vulcanFuncRegexp, vulcanFileRegexp, packRegexp, 5)
	if err != nil {
		return fmt.Errorf("Can't create logger: %s", err.Error())
	}
	log.SetLogger(logger)

	if err := srv.Start(); err != nil {
		return fmt.Errorf("service start failure: %s", err)
	}

	return nil
}

func main() {
	r, err := registry.GetRegistry()
	if err != nil {
		fmt.Printf("Service exited with error: %s\n", err)
		os.Exit(255)
	}

	if err := Run(r); err != nil {
		fmt.Printf("Service exited with error: %s\n", err)
		os.Exit(255)
	} else {
		fmt.Println("Service exited gracefully")
	}
}
