package registry

import (
	"github.com/FoxComm/vulcand/plugin"

	"github.com/FoxComm/vulcand/plugin/connlimit"

	"github.com/FoxComm/vulcand/plugin/ratelimit"

	"github.com/FoxComm/vulcand/plugin/rewrite"

	"github.com/FoxComm/vulcand/plugin/cbreaker"

	"github.com/FoxComm/vulcand/plugin/trace"

	"github.com/FoxComm/core_services/router/plugins/feature_validator"
	"github.com/FoxComm/core_services/router/plugins/health_check"
	"github.com/FoxComm/core_services/router/plugins/recover_middleware"
	"github.com/FoxComm/core_services/router/plugins/session"
	"github.com/FoxComm/core_services/router/plugins/site_activity"
)

func GetRegistry() (*plugin.Registry, error) {
	r := plugin.NewRegistry()

	specs := []*plugin.MiddlewareSpec{

		connlimit.GetSpec(),

		ratelimit.GetSpec(),

		rewrite.GetSpec(),

		cbreaker.GetSpec(),

		trace.GetSpec(),

		health_check.GetSpec(),
		feature_validator.GetSpec(),
		session.GetSpec(),
		site_activity.GetSpec(),
		recover_middleware.GetSpec(),
	}

	for _, spec := range specs {
		if err := r.AddSpec(spec); err != nil {
			return nil, err
		}
	}
	return r, nil
}
