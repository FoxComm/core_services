package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/qor/qor/resource"

	"github.com/FoxComm/core_services/catalog_cache/lib/qor/catalog"
	"github.com/FoxComm/libs/configs"
	"github.com/FoxComm/libs/logger"
)

const cacheDuration = 1 * time.Minute

var (
	cache *ccache.Cache
	env   string
)

func init() {
	cache = ccache.New(ccache.Configure())
	env = configs.Get("FC_ENV")
}

func toFilterQuery(key, value string) string {
	return fmt.Sprintf("%v:%v", key, value)
}

func Search(params url.Values, solrHost string) (APIResponse, error) {
	solrURL, _ := url.Parse(fmt.Sprintf(solrHost+"/solr/%v/select?wt=json", env))

	solrQuery := solrURL.Query()
	solrQuery.Add("fq", toFilterQuery("type", `Spree\:\:Product`))

	page, prePage := 1, 30
	if pageStr := params.Get("page"); pageStr != "" {
		if i, err := strconv.Atoi(pageStr); err == nil {
			page = i
		}
	}

	if perPageStr := params.Get("per_page"); perPageStr != "" {
		if i, err := strconv.Atoi(perPageStr); err == nil {
			prePage = i
		}
	}
	solrQuery.Set("start", strconv.Itoa((page-1)*prePage))
	solrQuery.Set("rows", strconv.Itoa(prePage))

	defaultState := "active"
	if state := params.Get("state"); state != "" {
		defaultState = state
	}
	solrQuery.Add("fq", toFilterQuery("state_ss", defaultState))

	if id := params.Get("id"); id != "" {
		solrQuery.Add("q", "*:*")
		solrQuery.Add("fq", toFilterQuery("id_is", id))
	}

	if ids := params.Get("ids"); ids != "" {
		solrQuery.Add("q", "*:*")
		fqValue := fmt.Sprintf("(%v)", strings.Replace(ids, ",", "-OR-", -1))
		solrQuery.Add("fq", toFilterQuery("id_is", fqValue))
	}

	if taxon_ids := params.Get("taxon_ids"); taxon_ids != "" {
		solrQuery.Add("q", "*:*")
		if query_type := params.Get("query_type"); query_type == "all_of" {
			for _, id := range strings.Split(taxon_ids, ",") {
				solrQuery.Add("fq", toFilterQuery("taxon_ids_ims", id))
			}
		} else {
			fqValue := fmt.Sprintf("(%v)", strings.Replace(taxon_ids, ",", "-OR-", -1))
			solrQuery.Add("fq", toFilterQuery("taxon_ids_ims", fqValue))
		}
	}

	if taxons, ok := params["taxons"]; ok {
		solrQuery.Add("q", "*:*")
		for _, taxon := range taxons {
			solrQuery.Add("fq", toFilterQuery("taxons_sms", fmt.Sprintf("(%v)", taxon)))
		}
	}

	if taxons, ok := params["permalink"]; ok {
		solrQuery.Add("q", "*:*")
		for _, taxon := range taxons {
			solrQuery.Add("fq", toFilterQuery("taxons_with_parent_sms", fmt.Sprintf("(%v)", taxon)))
		}
	}

	if categories, ok := params["categories"]; ok {
		solrQuery.Add("q", "*:*")
		fqValue := fmt.Sprintf("(%v)", strings.Join(categories, "-OR-"))
		solrQuery.Add("fq", toFilterQuery("taxons_sms", fqValue))
	}

	if slug := params.Get("slug"); slug != "" {
		solrQuery.Add("q", "*:*")
		solrQuery.Add("fq", toFilterQuery("slug_ss", slug))
	}

	if sort := params.Get("sort"); sort != "" {
		sortDirection := "asc"
		if params.Get("sort_desc") != "" {
			sortDirection = "desc"
		}
		solrQuery.Add("sort", sort+"JOIN"+sortDirection)
	}

	if keywords, ok := params["keywords"]; ok {
		solrQuery.Add("qf", "name_texts description_texts meta_description_texts")
		solrQuery.Add("fl", "* score")
		solrQuery.Add("defType", "edismax")
		for _, keyword := range keywords {
			solrQuery.Add("q", keyword)
		}
		solrURL.RawQuery = solrQuery.Encode()
	} else {
		// Solr needs like this: q:(escaped\+string+OR+another\+escaped\+string)
		if solrQuery.Get("q") == "" {
			solrQuery.Add("q", "*:*")
		}
		solrURL.RawQuery = strings.Replace(strings.Replace(solrQuery.Encode(), "+", "\\+", -1), "-OR-", "+OR+", -1)
	}

	// Ugly hack to avoid '+' scaping :(
	solrURL.RawQuery = strings.Replace(solrURL.RawQuery, "JOIN", "+", -1)

	url := solrURL.String()

	if item := cache.Get(url); item != nil && !item.Expired() {
		resp, ok := item.Value().(APIResponse)
		if ok {
			return resp, nil
		} else {
			logger.Error("[catalog_cache] Can't cast catalog_cache cache value to response")
			cache.Delete(url)
		}
	}

	if resp, err := http.Get(url); err == nil {
		var result Result
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		values := map[string]interface{}{}
		if err := decoder.Decode(&values); err == nil {
			resource.DecodeToResource(solr, &result, catalog.ConvertMapToMetaValues(values, solrQor), nil).Start()
			result.Response.CurrentPage = page
			result.Response.PerPage = prePage
			result.Response.Pages = (result.Response.TotalCount-1)/prePage + 1
			cache.Set(url, result.Response, cacheDuration)
			return result.Response, nil
		} else {
			return APIResponse{Results: []Product{}}, err
		}
	} else {
		return APIResponse{Results: []Product{}}, err
	}
}
