package test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/feature_manager/core"
	"github.com/FoxComm/FoxComm/router/registry"
	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/plugin"
	"github.com/FoxComm/vulcand/service"
)

var srvPort = 18080

func NewTestMemEngineServer(t *testing.T) *VulcanServer {
	srv, err := NewMemEngineServer()
	require.NoError(t, err)
	srv.T = t
	return srv
}

func NewMemEngineServer() (*VulcanServer, error) {
	var err error
	port := configs.GetSafeRunPortStringFromString(fmt.Sprintf("%d", srvPort))
	srvPort += 1

	options := service.Options{}
	options.EngineType = "mem"
	options.Log = "console"
	options.LogSeverity = "WARN"
	testRegistry, err := registry.GetRegistry()
	if err != nil {
		return nil, err
	}

	srv := service.NewService(options, testRegistry)

	return &VulcanServer{
		Registry: testRegistry,
		Service:  srv,
		Port:     port,
	}, err
}

type VulcanServer struct {
	Port     string
	Service  *service.Service
	Registry *plugin.Registry
	T        *testing.T

	changesC chan interface{}
	stopC    chan bool
}

func (srv *VulcanServer) Start() {
	go func() {
		err := srv.Service.Start()
		if srv.T != nil {
			assert.NoError(srv.T, err)
		}
	}()
	srv.Service.WaitUntilStarted()
	ng := srv.Service.GetEngine()
	if srv.T != nil {
		assert.NotNil(srv.T, ng, "Engine not nil")
	}

	// add listener
	l := engine.Listener{
		Id:       "default",
		Protocol: "http",
		Address: engine.Address{
			Address: fmt.Sprintf("0.0.0.0:%s", srv.Port),
			Network: "tcp",
		},
	}
	err := ng.UpsertListener(l)
	if srv.T != nil {
		assert.NoError(srv.T, err)
	}

	srv.changesC = make(chan interface{})
	srv.stopC = make(chan bool)
	//go ng.Subscribe(srv.changesC, srv.stopC)
}

func (srv *VulcanServer) Stop() {
	ng := srv.Service.GetEngine()
	ng.DeleteListener(engine.ListenerKey{Id: "default"})
	srv.Service.Stop()
}

func (srv *VulcanServer) BuildUrl(url string) string {
	return fmt.Sprintf("http://localhost:%s/%s", srv.Port, url)
}

func (srv *VulcanServer) Get(url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", srv.BuildUrl(url), nil)
	if srv.T != nil {
		require.NoError(srv.T, err)
	}
	req.Header.Add("Accept", "text/html")
	resp, err := client.Do(req)
	if srv.T != nil {
		require.NoError(srv.T, err)
	}
	return resp
}

func (srv *VulcanServer) UseRouterMiddleware(key engine.FrontendKey, mtype string, priority int) engine.Middleware {
	spec := srv.Registry.GetSpec(mtype)
	pl, err := spec.FromJSON([]byte("{}"))

	if srv.T != nil {
		require.NoError(srv.T, err, "Can't create plugin.Middleware")
	}

	m := CreateMiddleware(mtype, mtype, priority, pl)
	ng := srv.Service.GetEngine()
	err = ng.UpsertMiddleware(key, m, 0)
	if srv.T != nil {
		assert.NoError(srv.T, err, "Can't insert router middleware")
	}
	return m
}

func (srv *VulcanServer) MockStoreId(key engine.FrontendKey, id int) {
	store := core.NewStoreByID(id)
	var feature *core.StoreFeature
	store.LoadFeatures()
	for _, f := range store.StoreFeatures {
		if f.FeatureName == key.Id {
			feature = &f
			break
		}
	}

	r := RouterMiddleware{
		HandlerFn: func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("FC-Store-ID", strconv.Itoa(id))
			// r.Header.Set("FC-Store-Admin-Spree-Token", store.SpreeToken)
			r.Header.Set("FC-Store-Host", "http://"+r.Host)
			r.Header.Set("FC-Solr-Host", store.SolrHost)
			if feature != nil {
				r.Header.Set("FC-Data-Source", feature.Datasource)
			}
		},
	}
	m := r.GetEngineMiddleware("feature_validator", "feature_validator", 0)
	ng := srv.Service.GetEngine()
	ng.UpsertMiddleware(key, m, 0)
}
