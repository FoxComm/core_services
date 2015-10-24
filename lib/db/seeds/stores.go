package seeds

import (
	"fmt"

	"github.com/FoxComm/core_services/feature_manager/core"
	"github.com/FoxComm/libs/configs"
	. "github.com/FoxComm/libs/db/masterdb"
	"github.com/FoxComm/libs/db/seeds"
	"github.com/FoxComm/libs/endpoints"
)

func (sr seeds.SeedRunner) GenerateStores() {
	env := configs.Get("FC_ENV")
	if env == "development" || env == "test" {
		var store core.Store

		if Db().First(&store, "name = ?", "Test Store").RecordNotFound() {
			merchant := core.Merchant{Name: "Test Merchant", Stores: []core.Store{
				{
					Name:     "Test Store",
					SolrHost: "http://localhost:8982",
					Domains: []core.Domain{
						{Domain: "http://localhost:8000"},
					},
				},
			}}
			Db().Create(&merchant)
			store = merchant.Stores[0]
		}

		for _, endpoint := range endpoints.Endpoints {
			if endpoint.IsFeature {
				feature := core.Feature{}
				Db().Where(&core.Feature{Name: endpoint.Name}).Assign("Description", endpoint.Description).FirstOrCreate(&feature)
				Db().Where(&core.StoreFeature{FeatureId: feature.Id, StoreId: store.Id}).
					Attrs("enabled", true).FirstOrCreate(&core.StoreFeature{})
			}
		}
	} else {
		fmt.Println("You can only create stores in the development or test environments")
	}

