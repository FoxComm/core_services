package test

import (
	"net/http"

	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/plugin"

	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
)

type RouterMiddleware struct {
	next        http.Handler
	PostHandler http.HandlerFunc
	HandlerFn   http.HandlerFunc
}

func (m *RouterMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	m.next = next
	return m, nil
}

func (m *RouterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.HandlerFn(w, r)
	m.next.ServeHTTP(w, r)
	if m.PostHandler != nil {
		m.PostHandler(w, r)
	}
}

func (m *RouterMiddleware) GetEngineMiddleware(id string, mtype string, priority int) engine.Middleware {
	return CreateMiddleware(id, mtype, priority, m)
}

func CreateMiddleware(id string, mtype string, priority int, m plugin.Middleware) engine.Middleware {
	return engine.Middleware{
		Id:         id,
		Type:       mtype,
		Middleware: m,
		Priority:   priority,
	}
}

func GenerateMiddlewareSpec(pluginType string, middleware plugin.Middleware) *plugin.MiddlewareSpec {
	fromOther := func(plugin.Middleware) (plugin.Middleware, error) {
		return middleware, nil
	}
	fromCli := func(c *cli.Context) (plugin.Middleware, error) {
		return fromOther(middleware)
	}

	return &plugin.MiddlewareSpec{
		Type:      pluginType,
		FromOther: fromOther,
		FromCli:   fromCli,
		CliFlags:  []cli.Flag{},
	}
}
