package tebex

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Based on https://documenter.getpostman.com/view/10912536/2s9XxvTEmh

const DefaultBaseUrl = "https://headless.tebex.io"

var (
	ErrHeadlessWebstoreNotFound    = fmt.Errorf("webstore not found")
	ErrHeadlessBasketNotFound      = fmt.Errorf("basket not found")
	ErrHeadlessCreatorCodeNotFound = fmt.Errorf("creator code not found")
)

var DefaultHeadlessClient = NewHeadlessClient(DefaultBaseUrl)

type HeadlessClient struct {
	url, privateKey string

	httpClient *http.Client
}

type HeadlessClientParams struct {
	Url        string // Required
	PrivateKey string // Optional
	HttpClient *http.Client
}

func NewHeadlessClient(url string) *HeadlessClient {
	return NewHeadlessClientWithOptions(HeadlessClientParams{Url: url})
}

func NewHeadlessClientWithOptions(params HeadlessClientParams) *HeadlessClient {
	c := &HeadlessClient{
		url:        params.Url,
		privateKey: params.PrivateKey,
		httpClient: params.HttpClient,
	}
	if strings.HasSuffix(c.url, "/") {
		c.url = c.url[:len(c.url)-1]
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

// MARK: Basket

func (c *HeadlessClient) CreateBasket(ctx context.Context, webstoreId string, body HeadlessCreateBasketRequest) (*HeadlessBasket, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/baskets", c.url, webstoreId)
	bodyRaw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyRaw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if body.IPAddress != "" {
		if c.privateKey == "" {
			return nil, fmt.Errorf("private key must be provided when IPAddress is set")
		}

		authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", webstoreId, c.privateKey)))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authString))
	}

	res, err := do[HeadlessBasket](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// MARK: Basket Packages

func (c *HeadlessClient) BasketAddPackage(ctx context.Context, basketId string, body HeadlessBasketAddPackageRequest) (*HeadlessBasket, error) {
	url := fmt.Sprintf("%s/api/baskets/%s/packages", c.url, basketId)
	bodyRaw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyRaw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := do[HeadlessBasket](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// MARK: Creator Codes

func (c *HeadlessClient) BasketApplyCreatorCode(ctx context.Context, webstoreId, basketId, creatorCode string) error {
	url := fmt.Sprintf("%s/api/accounts/%s/baskets/%s/creator-codes", c.url, webstoreId, basketId)
	body := fmt.Sprintf(`{"code":"%s"}`, creatorCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if _, err = do[any](c.httpClient, req); err != nil {
		return err
	}

	// Don't care about the response unless it was non-200 which was handled in `do`.
	return nil
}

func (c *HeadlessClient) BasketRemoveCreatorCode(ctx context.Context, webstoreId, basketId string) error {
	url := fmt.Sprintf("%s/api/accounts/%s/baskets/%s/creator-codes/remove", c.url, webstoreId, basketId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if _, err = do[any](c.httpClient, req); err != nil {
		return err
	}

	// Don't care about the response unless it was non-200 which was handled in `do`.
	return nil
}

type headlessResponse[T any] struct {
	// Always updated to the real http status code if zero (ie on success)
	Status int `json:"status"`

	// Success fields
	Message string `json:"message"`
	Data    *T     `json:"data"`

	// Failure fields
	Title  string `json:"title"`
	Detail string `json:"detail"`

	// There are other fields here, but they are empty as far as I(matt) can tell
}

func (r *headlessResponse[T]) into() error {
	if r.Status < 400 {
		return nil
	}

	if err, ok := knownErrors[r.Detail]; ok {
		return err
	}

	if r.Detail != "" {
		return errors.New(r.Detail)
	}

	return fmt.Errorf("bad response: %d", r.Status)
}

func do[T any](httpClient *http.Client, req *http.Request) (*headlessResponse[T], error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res headlessResponse[T]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if res.Status == 0 {
		res.Status = resp.StatusCode
	}
	if err = res.into(); err != nil {
		return nil, err
	}

	return &res, nil
}

// kno wnErrors is a map of the plain text error messages that are known/handled explicitly.
// For some unknown reason the tebex error gives no useful info besides the user facing/plaintext message.
var knownErrors = map[string]error{
	"Invalid account identifier provided": ErrHeadlessWebstoreNotFound,
	"Invalid basket identifier provided":  ErrHeadlessBasketNotFound,
	"Invalid creator code provided":       ErrHeadlessCreatorCodeNotFound,
}
