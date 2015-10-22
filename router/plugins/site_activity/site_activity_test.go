package site_activity_test

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/FoxComm/FoxComm/models"
	"github.com/FoxComm/FoxComm/repositories"
	"github.com/FoxComm/FoxComm/router/common/test"
	"github.com/FoxComm/FoxComm/utils"
	"github.com/FoxComm/vulcand/engine/tomlng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const storeId = 1

func TestSessionTokenIsSet(t *testing.T) {
	srv := test.NewTestMemEngineServer(t)
	srv.Start()
	defer srv.Stop()
	ng := srv.Service.GetEngine()

	s := test.MakeVulcanEndpoint(&test.SimpleCodeHandler{Code: 200})
	e, err := test.UpsertServer(ng, test.EndpointSettings{Url: s.URL, Name: "origin_frontend", Route: "/"})
	assert.NoError(t, err)

	srv.MockStoreId(e.Frontend.GetKey(), storeId)

	// Load middlewares for origin frontend from router config
	ngtoml, err := tomlng.New(srv.Registry, tomlng.Options{
		MainConfigFilepath: test.RouterRoot() + "/config.toml",
	})
	require.NoError(t, err)

	routerMiddlewares, err := ngtoml.GetMiddlewares(e.Frontend.GetKey())
	require.NoError(t, err)
	for _, rm := range routerMiddlewares {
		if rm.Type != "feature_validator" {
			ng.UpsertMiddleware(e.Frontend.GetKey(), rm, 0)
		}
	}

	shareToken := "TestSessionTokenIsSet" + utils.GenerateToken(4)

	resp := srv.Get(fmt.Sprintf("signup?fc_ut=%s&fc_at=4", shareToken))
	assert.Equal(t, 200, resp.StatusCode)

	// test SessionToken is set
	repo, err := repositories.NewSiteActivityRepoWithStoreId(storeId)
	require.NoError(t, err)

	var res []models.SiteActivity
	repo.Collection.Find(bson.M{"sharertoken": shareToken}).Sort("-timestamp").Limit(1).All(&res)
	require.NotEmpty(t, res)
	assert.NotEmpty(t, res[0].SessionToken, "SessionToken")
	assert.Equal(t, res[0].SharerToken, shareToken)
}
