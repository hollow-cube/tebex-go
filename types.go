package tebex

import "time"

type PaymentStatusType int

const (
	PaymentStatus__0 PaymentStatusType = iota // Unknown, is there a 0 entry??
	PaymentStatusComplete
	PaymentStatusRefund
	PaymentStatus__3 // Unknown
	PaymentStatus__4 // Unknown
	PaymentStatusCancelled
)

type PaymentStatus struct {
	Id          PaymentStatusType `json:"id"`
	Description string            `json:"description"`
}

type PaymentSequence string

const (
	PaymentSequenceOneOff PaymentSequence = "oneoff"
	PaymentSequenceFirst  PaymentSequence = "first"
	//todo what else?
)

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentMethod struct {
	Name       string `json:"name"`
	Refundable bool   `json:"refundable"`
}

type Variable struct {
	Identifier string `json:"identifier"`
	Option     string `json:"option"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Customer struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Ip               string `json:"ip"`
	Username         User   `json:"username"`
	MarketingConsent bool   `json:"marketing_consent"`
	Country          string `json:"country"`
	PostalCode       string `json:"postal_code"`
}

type ProductPurchase struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Quantity  int         `json:"quantity"`
	BasePrice Price       `json:"base_price"`
	PaidPrice Price       `json:"paid_price"`
	Variables []*Variable `json:"variables"`
	ExpiresAt *time.Time  `json:"expires_at"`
	Custom    *string     `json:"custom"`
	Username  User        `json:"username"`
}

type DeclineReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Payment struct {
	TransactionId             string             `json:"transaction_id"`
	Status                    PaymentStatus      `json:"status"`
	PaymentSequence           PaymentSequence    `json:"payment_sequence"`
	CreatedAt                 time.Time          `json:"created_at"`
	Price                     Price              `json:"price"`
	PricePaid                 Price              `json:"price_paid"`
	PaymentMethod             PaymentMethod      `json:"payment_method"`
	Fees                      map[string]Price   `json:"fees"`
	Customer                  Customer           `json:"customer"`
	Products                  []*ProductPurchase `json:"products"`
	Coupons                   []interface{}      `json:"coupons"`                     // todo
	GiftCards                 []interface{}      `json:"gift_cards"`                  // todo
	RecurringPaymentReference *string            `json:"recurring_payment_reference"` // Only present for recurring payments
	DeclineReason             *DeclineReason     `json:"decline_reason"`              // Only present for payment.declined (I think)
}

type RecurringPayment struct {
	Reference      string        `json:"reference"`
	CreatedAt      time.Time     `json:"created_at"`
	NextPaymentAt  time.Time     `json:"next_payment_at"`
	Status         PaymentStatus `json:"status"`
	InitialPayment Payment       `json:"initial_payment"`
	LastPayment    Payment       `json:"last_payment"`
	FailCount      int           `json:"fail_count"`
	Price          Price         `json:"price"`
	CancelledAt    *time.Time    `json:"cancelled_at"`  // Only present if Status.Id == PaymentStatusCancelled
	CancelReason   *string       `json:"cancel_reason"` // Only present if Status.Id == PaymentStatusCancelled
}
