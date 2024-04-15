# tebex-go

[![license](https://img.shields.io/github/license/hollow-cube/tebex-go.svg)](LICENSE)

A small Tebex webhook & partial headless API implementation for Go. Includes all the parts to 
validate and parse a Tebex webhook message, as well as a partial client implementation of the
headless api for creating baskets/checkout links.

## Install

```bash
go get github.com/hollow-cube/tebex-go@latest
```

## Webhook Processing

```go
func Handle(r *http.Request, w http.ResponseWriter) error {
    // Validate the payload from the request
    // May want to set the last arg (checkIp) to false if using a proxy. If true, it will check if the request
    // came from one of the known Tebex IPs. These IPs are also available at tebex.TebexWebhookIpAddresses.
    body, err := tebex.ValidatePayload(r, []byte("my-secret"), true)
    if err != nil {
        return err		
    }
	
    // Parse the payload 
    event, err := tebex.ParseEvent(body)
    if err != nil {
        return err
    }
    
    // Handle the event
    switch subject := event.Subject.(type) {
    case *tebex.PaymentCompletedEvent: 
        // Handle payment completed event
    }
    
    return nil
}
```

## Headless API

The headless API is not completely supported, the supported endpoints can be found below.
More information about the headless api can be found on the [official documentation](https://docs.tebex.io/developers/headless-api/overview).

```go
# Obtain the default client (using the official endpoint & http.DefaultClient)
headless := tebex.DefaultHeadlessClient

# Create a new basket
basket, err := headless.CreateBasket(ctx, myWebstoreId, tebex.HeadlessCreateBasketRequest{Username: "notmattw"})

# Add a package to the basket
basket, err := headless.BasketAddPackage(ctx, basket.Ident, tebex.HeadlessBasketAddPackageRequest{PackageId: 789, Quantity: 1})

# Add a creator code to the basket
err := headless.BasketApplyCreatorCode(ctx, myWebstoreId, basket.Ident, "myCreatorCode")

# Remove any applied creator code from the basket
err := headless.BasketRemoveCreatorCode(ctx, myWebstoreId, basket.Ident)

# Get the checkout link for the basket (will not be present until at least one package is added)
checkoutUrl := basket.Links.Checkout
```

## Contributing

Contibutions in the form of bug reports (via issues) or pull requests are welcome.
Discussion of the project can be done in the [Hollow Cube Discord Server](https://discord.hollowcube.net).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
