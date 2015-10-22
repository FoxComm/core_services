package feature_validator

import (
	"net/http"

	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/FoxComm/logger"
	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FoxComm/vulcand/engine"
	"github.com/FoxComm/vulcand/plugin"
)

const Type = "feature_validator"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

type FeatureValidatorPlugin struct {
	Endpoint *endpoints.Endpoint `json:"-"`
}

// Returns vulcan library compatible middleware
func (h *FeatureValidatorPlugin) NewHandler(next http.Handler) (http.Handler, error) {
	return New(next, h.Endpoint), nil

}

func (h *FeatureValidatorPlugin) OnUpsertToFrontend(f engine.Frontend) {
	ep, err := endpoints.Find(f.GetId())
	if err != nil {
		logger.Error("Can't find endpoint %s", f.GetId())
	}
	h.Endpoint = ep
}

func NewFeatureValidatorPlugin() (*FeatureValidatorPlugin, error) {
	return &FeatureValidatorPlugin{}, nil
}

func (c *FeatureValidatorPlugin) String() string {
	return "FeatureValidatorPlugin"
}

func FromOther(h FeatureValidatorPlugin) (plugin.Middleware, error) {
	return NewFeatureValidatorPlugin()
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return FromOther(FeatureValidatorPlugin{})
}

func CliFlags() []cli.Flag {
	return []cli.Flag{}
}
