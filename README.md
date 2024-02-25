# tebex-go

[![license](https://img.shields.io/github/license/hollow-cube/tebex-go.svg)](LICENSE)

A small Tebex webhook implementation for Go. Includes all the parts to validate and parse a
Tebex webhook message.

## Install

```bash
go get github.com/hollow-cube/tebex-go@latest
```

## Usage

The library can be used with

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

