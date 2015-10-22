package core

import (
	"errors"
	"net/url"
	"sync"

	. "github.com/FoxComm/libs/db/masterdb"
)

var DomainStoreMap map[string]*Store //This is a convenience map, for the sake of speed.

type Merchant struct {
	Id          int
	Name        string
	Description string
	Stores      []Store
}

type Store struct {
	Id            int
	MerchantId    int
	Name          string
	Description   string
	SpreeToken    string
	SolrHost      string
	OriginHost    string
	Domains       []Domain
	StoreFeatures []StoreFeature
}

type Domain struct {
	Id      int
	StoreId int
	Domain  string
}

type Feature struct {
	Id          int
	Name        string
	Description string
	Enabled     bool `sql:"-"`
}

//Mapping table
type StoreFeature struct {
	Id          int
	FeatureId   int
	FeatureName string `sql:"-"`
	StoreId     int
	Enabled     bool
	Datasource  string
}

func (store *Store) LoadFeatures() {
	Db().Select("feature_id, features.name as feature_name, store_id, enabled, datasource").
		Joins("inner join features on features.id = store_features.feature_id").
		Where("store_features.store_id = ?", store.Id).Find(&store.StoreFeatures)
}

var cacheMaplock = new(sync.RWMutex)

func ClearCacheMap() {
	cacheMaplock.Lock()
	defer cacheMaplock.Unlock()

	DomainStoreMap = nil
}

func generateCacheMaps() {
	cacheMaplock.RUnlock()
	cacheMaplock.Lock()
	defer cacheMaplock.RLock()
	defer cacheMaplock.Unlock()

	var domainStoreMap = map[string]*Store{}
	var stores []Store
	var merchants []Merchant

	Db().Find(&merchants)
	for _, merchant := range merchants {
		Db().Find(&merchant.Stores, "merchant_id = ?", merchant.Id)
		stores = append(stores, merchant.Stores...)
	}

	for _, store := range stores {
		Db().Find(&store.Domains, "store_id = ?", store.Id)
		store.LoadFeatures()
		for _, domain := range store.Domains {
			if u, err := url.Parse(domain.Domain); err == nil {
				domainStoreMap[u.Host] = &store
			}
		}
	}

	DomainStoreMap = domainStoreMap
}

func NewStoreByID(id int) *Store {
	var store Store
	Db().Find(&store, id)
	return &store
}

func CheckStoreAndFeature(domain string, featureName string) (store *Store, feature *StoreFeature, err error) {
	cacheMaplock.RLock()
	defer cacheMaplock.RUnlock()

	if DomainStoreMap == nil {
		generateCacheMaps()
	}

	if store, ok := DomainStoreMap[domain]; ok {
		for _, feature := range store.StoreFeatures {
			if feature.FeatureName == featureName {
				return store, &feature, nil
			}
		}
		return store, nil, nil
	} else {
		return nil, nil, errors.New("can't find related store")
	}
}
