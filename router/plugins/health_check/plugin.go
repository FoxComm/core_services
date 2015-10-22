package health_check

import (
	"net/http"

	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/plugin"
)

const Type = "health_check"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

type HealthCheckPlugin struct {
	Engine engine.Engine `json:"-"`
}

// Returns vulcan library compatible middleware
func (h *HealthCheckPlugin) NewHandler(next http.Handler) (http.Handler, error) {
	return New(next, h.Engine), nil
}

func (h *HealthCheckPlugin) InitEngine(ng engine.Engine) {
	h.Engine = ng
}

func NewHealthCheckPlugin() (*HealthCheckPlugin, error) {
	return &HealthCheckPlugin{}, nil
}

func (c *HealthCheckPlugin) String() string {
	return "HealthCheckPlugin"
}

func FromOther(h HealthCheckPlugin) (plugin.Middleware, error) {
	return NewHealthCheckPlugin()
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return FromOther(HealthCheckPlugin{})
}

func CliFlags() []cli.Flag {
	return []cli.Flag{}
}
