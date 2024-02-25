package tebex

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	ValidationEventType                    EventType = "validation.webhook"
	PaymentCompletedEventType              EventType = "payment.completed"
	PaymentDeclinedEventType               EventType = "payment.declined"
	PaymentRefundedEventType               EventType = "payment.refunded"
	PaymentDisputeOpenedEventType          EventType = "payment.dispute.opened"
	PaymentDisputeWonEventType             EventType = "payment.dispute.won"
	PaymentDisputeLostEventType            EventType = "payment.dispute.lost"
	PaymentDisputeClosedEventType          EventType = "payment.dispute.closed"
	RecurringPaymentStartedEventType       EventType = "recurring-payment.started"
	RecurringPaymentRenewedEventType       EventType = "recurring-payment.renewed"
	RecurringPaymentEndedEventType         EventType = "recurring-payment.ended"
	RecurringPaymentStatusChangedEventType EventType = "recurring-payment.status-changed"
)

type eventInternal struct {
	Id      string          `json:"id"`
	Type    EventType       `json:"type"`
	Date    time.Time       `json:"date"`
	Subject json.RawMessage `json:"subject"`
}

// Event is the common wrapper around any Tebex event
type Event struct {
	Id      string      `json:"id"`
	Type    EventType   `json:"type"`
	Date    time.Time   `json:"date"`
	Subject interface{} `json:"subject"`
}

// ValidationEvent is sent by Tebex when the webhook is added to the project
// to verify that it expects to receive Tebex webhooks.
//
// type = `validation.webhook`
type ValidationEvent struct {
}

type PaymentCompletedEvent Payment

type PaymentDeclinedEvent Payment

type PaymentRefundedEvent Payment

type PaymentDisputeOpenedEvent Payment

type PaymentDisputeWonEvent Payment

type PaymentDisputeLostEvent Payment

type PaymentDisputeClosedEvent Payment

type RecurringPaymentStartedEvent RecurringPayment

type RecurringPaymentRenewedEvent RecurringPayment

type RecurringPaymentEndedEvent RecurringPayment

type RecurringPaymentStatusChangedEvent RecurringPayment
