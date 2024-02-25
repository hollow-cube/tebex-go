package tebex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

var (
	ErrInvalidIp          = errors.New("invalid remote ip address")
	ErrInvalidContentType = errors.New("invalid content type")
	ErrMissingSignature   = errors.New("missing signature")
	ErrInvalidSignature   = errors.New("invalid signature")

	TebexWebhookIpAddresses = []string{"18.209.80.3", "54.87.231.232"}
)

func ValidatePayload(r *http.Request, secret []byte, checkIp bool) ([]byte, error) {
	contentType := r.Header.Get("content-type")
	signature := r.Header.Get("x-signature")

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var remoteIp string
	if checkIp {
		remoteIp = strings.Split(r.RemoteAddr, ":")[0]
	}

	return body, ValidatePayloadRaw(contentType, body, signature, secret, remoteIp)
}

type TebexWebhookData struct {
	ContentType string
	Secret      []byte
	Signature   string
	Body        []byte
}

func ValidatePayloadRaw(contentType string, body []byte, signature string, secret []byte, remoteAddr string) error {
	// Tebex only sends requests from a small number of IP addresses, so we can check that here.
	// Though when behind a proxy you may not have IP forwarding, so it's optional.
	if remoteAddr != "" {
		if !slices.Contains(TebexWebhookIpAddresses, remoteAddr) {
			return fmt.Errorf("%w: %s", ErrInvalidIp, remoteAddr)
		}
	}

	// Check that the required headers are present
	if contentType != "application/json" {
		return ErrInvalidContentType
	}
	if len(signature) == 0 {
		return ErrMissingSignature
	}

	// Validate the signature (if a signature was provided)
	return validateSignature(body, signature, secret)
}

func ParseEvent(payload []byte) (*Event, error) {
	var event eventInternal
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	subject, err := newEventStruct(event.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}
	if err = json.Unmarshal(event.Subject, &subject); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subject: %w", err)
	}

	return &Event{
		Id:      event.Id,
		Type:    event.Type,
		Date:    event.Date,
		Subject: subject,
	}, nil
}

func validateSignature(body []byte, signature string, secret []byte) error {
	// If no secret is provided, don't validate anything
	if len(secret) == 0 {
		return nil
	}

	// sha256 the body
	bodyHasher := sha256.New()
	bodyHasher.Write(body)
	bodyHash := hex.EncodeToString(bodyHasher.Sum(nil))

	// hmac the body hash with the secret
	hmacHasher := hmac.New(sha256.New, secret)
	hmacHasher.Write([]byte(bodyHash))
	computedSig := hex.EncodeToString(hmacHasher.Sum(nil))

	// Test the computed signature against the provided signature
	if computedSig != signature {
		return ErrInvalidSignature
	}
	return nil
}

func newEventStruct(eventType EventType) (interface{}, error) {
	switch eventType {
	case ValidationEventType:
		return &ValidationEvent{}, nil
	case PaymentCompletedEventType:
		return &PaymentCompletedEvent{}, nil
	case PaymentDeclinedEventType:
		return &PaymentDeclinedEvent{}, nil
	case PaymentRefundedEventType:
		return &PaymentRefundedEvent{}, nil
	case PaymentDisputeOpenedEventType:
		return &PaymentDisputeOpenedEvent{}, nil
	case PaymentDisputeWonEventType:
		return &PaymentDisputeWonEvent{}, nil
	case PaymentDisputeLostEventType:
		return &PaymentDisputeLostEvent{}, nil
	case PaymentDisputeClosedEventType:
		return &PaymentDisputeClosedEvent{}, nil
	case RecurringPaymentStartedEventType:
		return &RecurringPaymentStartedEvent{}, nil
	case RecurringPaymentRenewedEventType:
		return &RecurringPaymentRenewedEvent{}, nil
	case RecurringPaymentEndedEventType:
		return &RecurringPaymentEndedEvent{}, nil
	case RecurringPaymentStatusChangedEventType:
		return &RecurringPaymentStatusChangedEvent{}, nil
	}
	return nil, fmt.Errorf("unknown event type: %s", eventType)
}
