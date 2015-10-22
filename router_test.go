package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/FoxComm/FoxComm/router/common/test"

	"github.com/FoxComm/FoxComm/router/netutils"
	"github.com/FoxComm/FoxComm/router/plugins/recover_middleware"
	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/plugin"
	"github.com/FoxComm/vulcand/plugin/cbreaker"
)

type commonHandler struct {
	HandlerFn http.HandlerFunc
}

func (h commonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandlerFn(w, r)
}

var srvPort = 18080

type PanicMiddleware struct {
	next http.Handler
}

func (m *PanicMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	m.next = next
	return m, nil
}

func (m *PanicMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start panic")
	panic("middleware paniced")
}

func NewCbreakerMiddleware(t *testing.T) plugin.Middleware {
	condition := "ResponseCodeRatio(500, 600, 0, 600) > 0.5"
	fallback := map[string]interface{}{}
	fallback["type"] = "response"
	action := map[string]interface{}{}
	action["body"] = "come back later!"
	action["statuscode"] = http.StatusTeapot
	fallback["action"] = action
	fallbackDuration := 1 * time.Second
	recoveryDuration := 1 * time.Second
	checkPeriod := 100 * time.Millisecond

	spec, err := cbreaker.NewSpec(condition, fallback, nil, nil, fallbackDuration, recoveryDuration, checkPeriod)
	assert.NoError(t, err)

	plugin, err := cbreaker.FromOther(*spec)
	assert.NoError(t, err)

	return plugin
}

func TestMiddlewarePanic(t *testing.T) {
	srv := test.NewTestMemEngineServer(t)
	srv.Start()
	defer srv.Stop()
	ng := srv.Service.GetEngine()

	e, err := test.UpsertServer(ng, test.EndpointSettings{Url: "http://localhost:0", Name: "panic"})
	assert.NoError(t, err)

	panicMiddleware := test.CreateMiddleware("panic1", "panic1", 1, &PanicMiddleware{})
	ng.UpsertMiddleware(e.Frontend.GetKey(), panicMiddleware, 0)

	// assert.Panics(t, func() {
	// 	http.Get("http://localhost:18080/panic/")
	// }, "Without recover it should panic")

	recoverMiddleware := test.CreateMiddleware("recover1", "recover1", 0, recover_middleware.New())
	ng.UpsertMiddleware(e.Frontend.GetKey(), recoverMiddleware, 0)

	resp := srv.Get("panic/")
	assert.NotNil(t, resp, "Response")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestEndpointBadGateway(t *testing.T) {
	srv := test.NewTestMemEngineServer(t)
	srv.Start()
	defer srv.Stop()

	s := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 502})
	url := s.URL

	ng := srv.Service.GetEngine()

	e, err := test.UpsertServer(ng, test.EndpointSettings{Url: url, Name: "badgateway", Route: "/bd"})
	assert.NoError(t, err)

	// prepare cbreaker middleware
	plugin := NewCbreakerMiddleware(t)
	mcbreaker := test.CreateMiddleware("cbreaker", "cbreaker", 0, plugin)
	ng.UpsertMiddleware(e.Frontend.GetKey(), mcbreaker, 0)

	resp := srv.Get("bd")
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func TestNoServers(t *testing.T) {
	srv := test.NewTestMemEngineServer(t)
	srv.Registry.SetNoServersErrorHandler(&NoSrvHandler{})

	srv.Start()
	defer srv.Stop()

	ng := srv.Service.GetEngine()

	e, err := test.UpsertServer(ng, test.EndpointSettings{Url: "", Name: "noservers"})
	assert.NoError(t, err)

	resp := srv.Get("noservers")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	s := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 200})

	ng.UpsertServer(e.Backend.GetUniqueId(), s, 0)
	resp = srv.Get("noservers")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServersWhenOnlyOneIsGood(t *testing.T) {
	srv := test.NewTestMemEngineServer(t)
	srv.Start()
	defer srv.Stop()

	ng := srv.Service.GetEngine()

	s1 := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 502})
	e, err := test.UpsertServer(ng, test.EndpointSettings{Url: s1.URL, Name: "twoservers", Route: "/ep"})
	assert.NoError(t, err)

	s2 := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 200})
	err = ng.UpsertServer(e.Backend.GetUniqueId(), s2, 0)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		resp := srv.Get("ep")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestHostHeaderIsPassed(t *testing.T) {
	// create vulcand instance
	srv := test.NewTestMemEngineServer(t)
	srv.Start()
	ng := srv.Service.GetEngine()
	defer srv.Stop()

	// Create and upsert server
	s := test.MakeVulcanEndpoint(commonHandler{HandlerFn: func(w http.ResponseWriter, r *http.Request) {
		netutils.WriteTextResponse(w, 200, r.Header.Get("FC-Host"))
	}})

	e, err := test.UpsertServer(ng, test.EndpointSettings{
		Url:      s.URL,
		Name:     "test",
		Frontend: &engine.HTTPFrontendSettings{PassHostHeader: true},
	})
	assert.NoError(t, err)

	// Create and upsert middleware to server
	m := test.RouterMiddleware{HandlerFn: func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("FC-Host", r.Host)
	}}
	vm := m.GetEngineMiddleware("f1", "f1", 1)
	ng.UpsertMiddleware(e.Frontend.GetKey(), vm, 0)

	// make requests
	client := http.Client{}

	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", srv.BuildUrl("test"), nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	// check all is ok
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, req.Host, string(body))
}

// make endpoints benchmarks?
func BenchmarkSelf(b *testing.B) {
	srv, _ := test.NewMemEngineServer()
	srv.Start()
	defer srv.Stop()

	ng := srv.Service.GetEngine()

	s1 := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 502})
	test.UpsertServer(ng, test.EndpointSettings{Url: s1.URL, Name: "bench"})

	for i := 0; i < b.N; i++ {
		srv.Get("bench")
	}
}
