package models

import (
	"time"

	"github.com/FoxComm/libs/spree"
)

//I didn't want to actually store all the fields, but I ended up storing so much, that I thought to include everything.
type CheckoutActivityDetails struct {
	Id                        int
	Number                    string  // customer-facing order number
	ItemTotal                 float64 `json:"item_total,string"` // this is not qty, that's below
	Total                     float64 `json:",string"`
	ShipTotal                 float64 `json:"ship_total,string"`
	State                     string
	AdjustmentTotal           float64 `json:"adjustment_total,string"`
	User                      spree.User
	CreatedAt                 time.Time `json:"created_at"`
	UpdateAt                  time.Time `json:"updated_at"`
	CompletedAt               time.Time `json:"completed_at"`
	PaymentTotal              float64   `json:"payment_total,string"`
	ShipmentState             string    `json:"shipment_state"`
	PaymentState              string    `json:"payment_state"`
	Email                     string
	SpecialInstructions       string `json:"special_instructions"`
	Channel                   string
	IncludedTaxTotal          float64 `json:"included_tax_total,string"`
	AdditionalTaxTotal        float64 `json:"additional_tax_total,string"`
	DisplayIncludedTaxTotal   string  `json:"display_included_tax_total"`
	DisplayAdditionalTaxTotal string  `json:"display_additional_tax_total"`
	TaxTotal                  float64 `json:"tax_total,string"`
	Currency                  string
	TotalQuantity             int    `json:"total_quantity"`
	DisplayTotal              string `json:"display_total"`
	DisplayShipTotal          string `json:"display_ship_total"`
	DisplayTaxTotal           string `json:"display_tax_total"`
	Token                     string
	BillAddress               spree.Address   `json:"bill_address"`
	ShipAddress               spree.Address   `json:"ship_address"`
	LineItems                 []OrderLineItem `json:"line_items"`
	Payments                  []OrderPayment
	Shipments                 []OrderShipment
	Adjustments               []OrderAdjustment
}

type OrderLineItem struct {
	OrderId   int `json:"order_id"`
	VariantId int `json:"variant_id"`
	Quantity  int
	Price     float64 `json:",string"`
}

type OrderPayment struct {
	Amount        float64   `json:",string"`
	CreatedAt     time.Time `json:"created_at"`
	Id            int
	PaymentMethod PaymentMethod `json:"payment_method"`
	State         string
}

type OrderShipment struct {
	//Someone fill this in
}

type OrderAdjustment struct {
	//Someone fill this in
}

type PaymentMethod struct {
	Environment string
	Id          int
	Name        string
}
