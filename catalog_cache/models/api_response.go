package models

import "encoding/json"

type Result struct {
	Response APIResponse
}

type APIResponse struct {
	Count       int       `json:"count"`
	TotalCount  int       `json:"total_count"`
	CurrentPage int       `json:"current_page"`
	PerPage     int       `json:"per_page"`
	Pages       int       `json:"pages"`
	Results     []Product `json:"results"`
}

type Product struct {
	Id                 int               `json:"id"`
	Name               string            `json:"name"`
	Brand              string            `json:"brand"`
	BrandPermalink     string            `json:"brand_permalink"`
	Price              float64           `json:"price"`
	Taxons             []string          `json:"taxons"`
	TaxonIds           []int             `json:"taxon_ids"`
	SimilarProductsIds json.RawMessage   `json:"similar_products_ids"`
	Description        string            `json:"description"`
	Slug               string            `json:"slug"`
	Master             Variant           `json:"master"`
	Variants           []*Variant        `json:"variants"`
	ProductProperties  []ProductProperty `json:"product_properties"`
	OptionTypes        []OptionType      `json:"option_types"`
	Rating             int               `json:"rating"`
}

type Variant struct {
	Id           int             `json:"id"`
	IsMaster     bool            `json:"is_master"`
	OptionsText  string          `json:"options_text"`
	Sku          string          `json:"sku"`
	Position     int             `json:"position"`
	Price        string          `json:"price"`
	Slug         string          `json:"slug"`
	InStock      bool            `json:"in_stock"`
	Images       json.RawMessage `json:"images"`
	OptionValues json.RawMessage `json:"option_values"`
}

type ProductProperty struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"$type"`
}

type OptionType struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	Presentation string          `json:"presentation"`
	Position     int             `json:"position"`
	OptionValues json.RawMessage `json:"option_values"`
}
