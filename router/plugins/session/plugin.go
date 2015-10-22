package session

import (
	"net/http"

	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FoxComm/vulcand/plugin"
)

const Type = "session"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

type SessionPlugin struct {
}

// Returns vulcan library compatible middleware
func (h *SessionPlugin) NewHandler(next http.Handler) (http.Handler, error) {
	return New(next), nil
}

func NewSessionPlugin() (*SessionPlugin, error) {
	return &SessionPlugin{}, nil
}

func (c *SessionPlugin) String() string {
	return "SessionPlugin"
}

func FromOther(h SessionPlugin) (plugin.Middleware, error) {
	return NewSessionPlugin()
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return FromOther(SessionPlugin{})
}

func CliFlags() []cli.Flag {
	return []cli.Flag{}
}
