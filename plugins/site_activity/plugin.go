package site_activity

import (
	"net/http"

	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FoxComm/vulcand/plugin"
)

const Type = "site_activity"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

type SiteActivityPlugin struct {
}

// Returns vulcan library compatible middleware
func (h *SiteActivityPlugin) NewHandler(next http.Handler) (http.Handler, error) {
	return New(next), nil
}

func NewSiteActivityPlugin() (*SiteActivityPlugin, error) {
	return &SiteActivityPlugin{}, nil
}

func (c *SiteActivityPlugin) String() string {
	return "SiteActivityPlugin"
}

func FromOther(h SiteActivityPlugin) (plugin.Middleware, error) {
	return NewSiteActivityPlugin()
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return FromOther(SiteActivityPlugin{})
}

func CliFlags() []cli.Flag {
	return []cli.Flag{}
}
