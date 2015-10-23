package models

import (
	"encoding/json"

	"github.com/FoxComm/core_services/catalog_cache/lib/qor/catalog"
	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/resource"

	"regexp"
	"strings"
)

var solr *catalog.Resource
var solrQor *admin.Resource

func init() {
	solr = catalog.NewResource(Result{})
	api := catalog.NewResource(APIResponse{})
	product := catalog.NewResource(Product{})

	solr.RegisterMeta(resource.Meta{Name: "response", Resource: api})
	solr.AddProcessor(func(res interface{}, values *resource.MetaValues, context *qor.Context) error {
		result := res.(*Result)
		result.Response.Count = len(result.Response.Results)
		return nil
	})

	api.RegisterMeta(resource.Meta{Name: "numFound", Alias: "TotalCount"})
	api.RegisterMeta(resource.Meta{Name: "docs", Alias: "Results", Resource: product})
	normalProductAttrs := []string{
		"id_is", "taxon_ids_ims", "price_es", "slug_ss", "taxons_sms",
		"brand_ss", "brand_permalink_ss", "name_texts", "description_texts", "similar_products_ids_ss",
		"rating_is",
	}

	attrRegexp := regexp.MustCompile(`^(\w+)_[^_]+$`)
	for _, attr := range normalProductAttrs {
		column := attrRegexp.FindStringSubmatch(attr)[1] // column is id, taxon_ids...
		product.RegisterMeta(resource.Meta{Name: attr, Alias: column})
	}

	product.RegisterMeta(resource.Meta{
		Name: "master_ss",
		Setter: func(res interface{}, values *resource.MetaValues, context *qor.Context) {
			if product, ok := res.(*Product); ok {
				decoder := json.NewDecoder(strings.NewReader(values.Get("master_ss").Value.(string)))
				decoder.Decode(&product.Master)
			}
		}})

	product.RegisterMeta(resource.Meta{
		Name: "variants_sms",
		Setter: func(res interface{}, values *resource.MetaValues, context *qor.Context) {
			if product, ok := res.(*Product); ok {
				variants := []*Variant{}
				for _, encodedVariant := range values.Get("variants_sms").Value.([]interface{}) {
					variant := Variant{}
					decoder := json.NewDecoder(strings.NewReader(encodedVariant.(string)))
					decoder.Decode(&variant)
					variants = append(variants, &variant)
				}
				product.Variants = variants
			}
		}})

	product.RegisterMeta(resource.Meta{
		Name: "product_properties_ss",
		Setter: func(res interface{}, values *resource.MetaValues, context *qor.Context) {
			if product, ok := res.(*Product); ok {
				var properties []ProductProperty
				decoder := json.NewDecoder(strings.NewReader(values.Get("product_properties_ss").Value.(string)))
				decoder.Decode(&properties)
				for index := range properties {
					properties[index].Type = "ProductProperty:#Foundationall"
				}
				product.ProductProperties = properties
			}
		}})

	product.RegisterMeta(resource.Meta{
		Name: "option_types_ss",
		Setter: func(res interface{}, values *resource.MetaValues, context *qor.Context) {
			if product, ok := res.(*Product); ok {
				decoder := json.NewDecoder(strings.NewReader(values.Get("option_types_ss").Value.(string)))
				decoder.Decode(&product.OptionTypes)
			}
		}})

	solrQor = (*admin.Resource)(solr)
}
