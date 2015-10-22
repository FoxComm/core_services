package health_check

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FoxComm/FoxComm/router/common/test"
	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/engine/memng"
	"github.com/FoxComm/vulcand/plugin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareEngine1(ng *memng.Mem, t *testing.T) {
	_, err := test.UpsertServer(ng, "core", "", "http://localhost")
	assert.NoError(t, err)
}

func TestStatusFor(t *testing.T) {
	ng := memng.New(plugin.NewRegistry())
	prepareEngine1(ng.(*memng.Mem), t)
	hc := HealthCheck{Engine: ng}

	request, _ := http.NewRequest("GET", "/", nil)
	backends, err := ng.GetBackends()
	assert.NoError(t, err)
	code, message := hc.statusFor(request, backends)

	assert.Equal(t, http.StatusOK, code)

	hosts := []host{host{Okay: false, Host: "http://localhost", StatusCode: 0}}
	expected := status{AllOkay: false, Endpoints: map[string][]host{"core": hosts}}

	exp_json, err := json.Marshal(expected)
	assert.NoError(t, err)

	msg_json, err := json.Marshal(message)
	assert.NoError(t, err)

	assert.Equal(t, string(exp_json), string(msg_json))
}

func TestRegister(t *testing.T) {
	router := gin.New()

	Register(router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health_check", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestIsServerGood(t *testing.T) {
	router := gin.New()
	Register(router)
	ts := httptest.NewServer(router)

	srv := engine.Server{Id: "srv2", URL: ts.URL}
	ok, code := isServerGood(srv)

	assert.True(t, ok)
	assert.Equal(t, 200, code)
}
