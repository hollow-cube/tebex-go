package tebex

type HeadlessCreateBasketRequest struct {
	CompleteUrl          string         `json:"complete_url,omitempty"`
	CompleteAutoRedirect bool           `json:"complete_auto_redirect,omitempty"`
	CancelUrl            string         `json:"cancel_url,omitempty"`
	Custom               map[string]any `json:"custom,omitempty"`

	// Should be present for Minecraft webstores.
	Username string `json:"username,omitempty"`
}

type HeadlessBasketAddPackageRequest struct {
	PackageId    int            `json:"package_id"`
	Quantity     int            `json:"quantity"`
	VariableData map[string]any `json:"variable_data,omitempty"`
}

type HeadlessBasket struct {
	Ident       string `json:"ident"`
	Complete    bool   `json:"complete"`
	Id          int    `json:"id"`
	Country     string `json:"country"`
	Ip          string `json:"ip"`
	UsernameId  string `json:"username_id"`
	Username    string `json:"username"`
	BasePrice   int    `json:"base_price"`
	SalesTax    int    `json:"sales_tax"`
	TotalPrice  int    `json:"total"`
	Packages    []any  `json:"packages"`
	Coupons     []any  `json:"coupons"`
	GiftCards   []any  `json:"gift_cards"`
	CreatorCode string `json:"creator_code"`
	Links       struct {
		Checkout string `json:"checkout"`
	} `json:"links"`
}
