package catalog

import (
	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/resource"
)

type Resource admin.Resource

func NewResource(value interface{}) *Resource {
	res := Resource(*admin.NewResource(value))
	return &res
}

func (res *Resource) AllAttrs() []*resource.Meta {
	r := admin.Resource(*res)
	return r.AllAttrs()
}

func (res *Resource) RegisterMeta(meta resource.Meta) {
	meta.Type = "-"
	meta.Value = func(interface{}, *qor.Context) interface{} { return nil }
	res.Resource.RegisterMeta(&meta)
}

func ConvertMapToMetaValues(values map[string]interface{}, res *admin.Resource) *resource.MetaValues {
	return admin.ConvertMapToMetaValues(values, res)
}
