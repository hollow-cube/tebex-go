package tebex

import "encoding/json"

type HeadlessCreateBasketRequest struct {
	CompleteUrl          string         `json:"complete_url,omitempty"`
	CompleteAutoRedirect bool           `json:"complete_auto_redirect,omitempty"`
	CancelUrl            string         `json:"cancel_url,omitempty"`
	Custom               map[string]any `json:"custom,omitempty"`

	// Username should be present for Minecraft webstores.
	Username string `json:"username,omitempty"`
	// IPAddress should be present if creating the request from a backend server
	IPAddress string `json:"ip_address,omitempty"`
}

type HeadlessBasketAddPackageRequest struct {
	PackageId    int            `json:"package_id"`
	Quantity     int            `json:"quantity"`
	VariableData map[string]any `json:"variable_data,omitempty"`
}

type HeadlessBasket struct {
	Ident       string               `json:"ident"`
	Complete    bool                 `json:"complete"`
	Id          int                  `json:"id"`
	Country     string               `json:"country"`
	Ip          string               `json:"ip"`
	UsernameId  string               `json:"username_id"`
	Username    string               `json:"username"`
	BasePrice   float64              `json:"base_price"`
	SalesTax    float64              `json:"sales_tax"`
	TotalPrice  float64              `json:"total"`
	Packages    []any                `json:"packages"`
	Coupons     []any                `json:"coupons"`
	GiftCards   []any                `json:"gift_cards"`
	CreatorCode string               `json:"creator_code"`
	Links       *HeadlessBasketLinks `json:"links"`
}

type HeadlessBasketLinks struct {
	Checkout string `json:"checkout"`
}

func (h *HeadlessBasketLinks) UnmarshalJSON(raw []byte) error {
	// For some reason they return an empty array when this is empty, so handle that case.
	if string(raw) == "[]" {
		return nil
	}

	type Alias HeadlessBasketLinks
	return json.Unmarshal(raw, (*Alias)(h))
}
