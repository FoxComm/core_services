package health_check

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/FoxComm/FoxComm/router/netutils"

	_ "github.com/FoxComm/FoxComm/utils/ssl"
	"github.com/FoxComm/vulcand/engine"
	"github.com/gin-gonic/gin"
)

type host struct {
	Host       string
	Okay       bool
	StatusCode int
}

type status struct {
	AllOkay   bool
	Endpoints map[string][]host
}

func newStatus() *status {
	s := &status{}
	s.Endpoints = map[string][]host{}
	return s
}

// Health Check Middleware
var timeout = time.Duration(1 * time.Second)
var httpTransport = http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		// We tell Dialer to reuse the same connection for 5 seconds
		d := net.Dialer{Timeout: timeout, KeepAlive: 5 * time.Second}
		return d.Dial(network, addr)
	},
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var httpClient = http.Client{Transport: &httpTransport}

type HealthCheck struct {
	next   http.Handler
	Engine engine.Engine
}

func New(next http.Handler, ng engine.Engine) *HealthCheck {
	h := &HealthCheck{next: next, Engine: ng}
	return h
}

func writeErr(w http.ResponseWriter, err error) {
	netutils.WriteJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func (h *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

	// self health check
	if len(paths) == 1 {
		s := newStatus()
		s.AllOkay = true
		netutils.WriteJsonResponse(w, http.StatusOK, s)
		return
	}

	if len(paths) >= 2 {
		featureName := paths[1]
		var backends []engine.Backend
		var err error

		switch featureName {
		case "all":
			backends, err = h.Engine.GetBackends()
			if err != nil {
				writeErr(w, err)
				return
			}
		default:
			frontend, err := h.Engine.GetFrontend(engine.FrontendKey{Id: featureName})
			if err != nil {
				writeErr(w, err)
				return
			}

			backend, err := h.Engine.GetBackend(engine.BackendKey{Id: frontend.BackendId})
			if err != nil {
				writeErr(w, err)
				return
			}
			backends = []engine.Backend{*backend}
		}

		code, message := h.statusFor(r, backends)
		netutils.WriteJsonResponse(w, code, message)
		return
	}
	netutils.WriteJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
	return
}

// statusFor checks status of all registered backends and outputs their status
func (h HealthCheck) statusFor(r *http.Request, backends []engine.Backend) (int, interface{}) {
	response := newStatus()
	response.AllOkay = true

	for _, backend := range backends {
		servers, err := h.Engine.GetServers(backend.GetUniqueId())
		if err != nil {
			return http.StatusInternalServerError, err
		}

		for _, server := range servers {
			ok, code := isServerGood(server)
			if !ok {
				response.AllOkay = false
			}
			h := host{Host: server.URL, Okay: ok, StatusCode: code}
			response.Endpoints[backend.Id] = append(response.Endpoints[server.URL], h)
		}
	}
	return http.StatusOK, response
}

func Register(engine *gin.Engine) {
	engine.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func isServerGood(server engine.Server) (bool, int) {
	ep, err := url.Parse(server.URL)
	if err != nil {
		return false, 0
	}

	ep.Path = "/health_check"
	//logger.Debug("[Health Check] Checking status of %s", url.String())

	resp, err := httpClient.Get(ep.String())

	if resp != nil {
		defer resp.Body.Close()
		return err == nil, resp.StatusCode
	}

	return false, 0
}
