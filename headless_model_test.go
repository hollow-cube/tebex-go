package tebex

import (
	"encoding/json"
	"testing"
)

func TestHeadlessBasketUnmarshalLinks1(t *testing.T) {
	data := `{
    "ident": "abcdef",
    "complete": true,
    "id": 12345,
    "country": "US",
    "ip": "127.0.0.1",
    "username_id": "abcdef",
    "username": "test123",
    "cancel_url": "https:\/\/hollowcube.net",
    "complete_url": "https:\/\/hollowcube.net",
    "complete_auto_redirect": true,
    "base_price": 39.96,
    "sales_tax": 0,
    "total_price": 39.96,
    "currency": "USD",
    "packages": [],
    "coupons": [],
    "giftcards": [],
    "creator_code": null,
    "links": {
      "checkout": "IAMLINK"
    }
  }`
	var result HeadlessBasket
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if result.Links.Checkout != "IAMLINK" {
		t.Errorf("Unexpected value: %v", result.Links.Checkout)
	}
}

func TestHeadlessBasketUnmarshalLinks2(t *testing.T) {
	data := `{
    "ident": "abcdef",
    "complete": true,
    "id": 12345,
    "country": "US",
    "ip": "127.0.0.1",
    "username_id": "abcdef",
    "username": "test123",
    "cancel_url": "https:\/\/hollowcube.net",
    "complete_url": "https:\/\/hollowcube.net",
    "complete_auto_redirect": true,
    "base_price": 39.96,
    "sales_tax": 0,
    "total_price": 39.96,
    "currency": "USD",
    "packages": [],
    "coupons": [],
    "giftcards": [],
    "creator_code": null,
    "links": []
  }`
	var result HeadlessBasket
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if result.Links.Checkout != "" {
		t.Errorf("Unexpected value: %v", result.Links.Checkout)
	}
}
