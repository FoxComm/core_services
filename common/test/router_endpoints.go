package test

import (
	"fmt"
	"net/http"

	"github.com/FoxComm/vulcand/engine"
)

func newFrontend(name, prefix string, settings *engine.HTTPFrontendSettings) (*engine.Frontend, error) {
	route := "PathRegexp(`^" + prefix + ".*`)"
	if settings == nil {
		settings = &engine.HTTPFrontendSettings{}
	}
	return engine.NewHTTPFrontend(name, name, route, *settings)
}

func newBackend(name string, settings *engine.HTTPBackendSettings) (*engine.Backend, error) {
	if settings == nil {
		settings = &engine.HTTPBackendSettings{}
	}
	return engine.NewHTTPBackend(name, *settings)
}

type Endpoint struct {
	Frontend engine.Frontend
	Backend  engine.Backend
	Server   engine.Server
}

type EndpointSettings struct {
	Url      string
	Name     string
	Route    string
	Frontend *engine.HTTPFrontendSettings
	Backend  *engine.HTTPBackendSettings
}

func UpsertServer(ng engine.Engine, settings EndpointSettings) (*Endpoint, error) {
	if settings.Route == "" {
		settings.Route = "/" + settings.Name
	}

	b, err := newBackend(settings.Name, settings.Backend)
	if err != nil {
		return nil, err
	}
	if err := ng.UpsertBackend(*b); err != nil {
		return nil, err
	}

	f, err := newFrontend(settings.Name, settings.Route, settings.Frontend)
	if err != nil {
		return nil, err
	}
	if err := ng.UpsertFrontend(*f, 0); err != nil {
		return nil, err
	}

	var s engine.Server
	if settings.Url != "" {
		s = engine.Server{Id: fmt.Sprintf("%ssrv1", settings.Name), URL: settings.Url}
		if err := ng.UpsertServer(b.GetUniqueId(), s, 0); err != nil {
			return nil, err
		}
	}
	return &Endpoint{Frontend: *f, Backend: *b, Server: s}, nil
}

func MakeVulcanEndpoint(handler http.Handler) engine.Server {
	s := MakeHttpServer(handler)

	return engine.Server{
		Id:  s.Addr[1:],
		URL: fmt.Sprintf("http://localhost%s", s.Addr),
	}
}
