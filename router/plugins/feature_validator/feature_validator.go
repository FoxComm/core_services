package feature_validator

import (
	"net/http"
	"strconv"

	"github.com/FoxComm/libs/endpoints"
	"github.com/FoxComm/core_services/feature_manager/core"
)

type FeatureValidator struct {
	Endpoint *endpoints.Endpoint
	next     http.Handler
}

func New(next http.Handler, endpoint *endpoints.Endpoint) *FeatureValidator {
	return &FeatureValidator{next: next, Endpoint: endpoint}
}

func (fv *FeatureValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if store, feature, err := core.CheckStoreAndFeature(r.Host, fv.Endpoint.Name); err == nil {
		if fv.Endpoint.IsFeature && !feature.Enabled {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
		}

		r.Header.Set("FC-Store-ID", strconv.Itoa(store.Id))
		r.Header.Set("FC-Store-Admin-Spree-Token", store.SpreeToken)

		scheme := "http"
		if r.Header.Get("X-Forwarded-Port") == "443" {
			scheme = "https"
		}

		r.Header.Set("FC-Store-Host", scheme+"://"+r.Host)
		r.Header.Set("FC-Solr-Host", store.SolrHost)
		if feature != nil {
			r.Header.Set("FC-Data-Source", feature.Datasource)
		}

		fv.next.ServeHTTP(w, r)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Store Not Found"))
	}
}
