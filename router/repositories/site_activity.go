package repositories

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/libs/db/db_switcher"
)

var SiteActivityCollection = "site_activity"

type SiteActivityRepo struct {
	db_switcher.Mongo
}

func NewSiteActivityRepo(request *http.Request) (*SiteActivityRepo, error) {
	repo := SiteActivityRepo{}

	err := repo.InitializeWithRequest(request, SiteActivityCollection)
	return &repo, err
}

func NewSiteActivityRepoWithStoreId(storeID int) (*SiteActivityRepo, error) {
	var repo SiteActivityRepo
	err := repo.InitializeWithStoreID(storeID, SiteActivityCollection, endpoints.SocialAnalyticsAPI)
	return &repo, err
}

func (repo *SiteActivityRepo) FindInboundActivities(sessionToken string, result interface{}) error {
	query := bson.M{
		"apirequesturl": "",
		"sessiontoken":  sessionToken,
	}
	return repo.FindAll(query, result)
}
