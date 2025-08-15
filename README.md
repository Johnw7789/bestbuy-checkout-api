# Best Buy Checkout API

This is a feature complete login and checkout API client written in Go. Please note that in order to use this project, you must provide an Akamai API implementation, using the included Akamai adapter in the akamai package. See more details below in the overview.

## Quick Start

### Prerequisites

- Go 1.23.6 or later
- Valid Best Buy account credentials
- Akamai API key for bypassing bot protection

### Installation

1. Get the repository:
```bash
go get github.com/Johnw7789/bestbuy-checkout-api
```

2. Install dependencies:
```bash
go mod tidy
```

### Basic Usage

The full checkout flow follows a three step process:

1. **Cart Management** - Add items to cart and handle queue system, if a queue is present
2. **Authentication** - Login with MFA, for better security as well as for bypassing the queue's email code verification
3. **Checkout** - Complete the purchase flow, wait for 3DS redirect

Here's an example showing the full checkout flow:

```go
// Define status update callback
updateStatus := func(status string) {
    log.Println(status)
}

// Step 1: Configure checkout options
checkoutOpts := checkout.CheckoutOpts{
    SkuId:        "6614325", // Product SKU ID
    AkamaiApiKey: "your-akamai-api-key",
    UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
    Proxy:        "", // Optional proxy
        
    // Account information
    Email:       "test@gmail.com",
    PhoneNumber: "1234567890",

    // Shipping address
    ShippingFirstName: "John",
    ShippingLastName:  "Doe",
    ShippingAddress1:  "123 Main St",
    ShippingCity:      "City",
    ShippingStateCode: "State",
    ShippingZipCode:   "12345",

    // Billing address
    BillingFirstName: "John",
    BillingLastName:  "Doe",
    BillingAddress1:  "123 Main St",
    BillingCity:      "City",
    BillingStateCode: "State",
    BillingZipCode:   "12345",

    // Payment information
    CardNumber:   "4242424242424242",
    CardExpMonth: "12",
    CardExpYear:  "25",
    CVV:          "123",
}

// Step 2: Create Akamai adapter
akamaiAdapter := shr.NewAkamaiAdapterWithHyper(checkoutOpts.AkamaiApiKey)

// Step 3: Create checkout client
checkoutClient, err := checkout.NewCheckoutClient(checkoutOpts, nil, akamaiAdapter, updateStatus)
if err != nil {
    log.Fatal("Failed to create checkout client:", err)
}

// Step 4: Add item to cart and handle queue if one is present
tokenId, a2cId, checkoutData, err := checkoutClient.CartItem()
if err != nil && err.Error() != "no queue found" {
    log.Fatal("Failed to complete cart process:", err)
}

// Step 5: Configure login options
loginOpts := login.LoginOpts{
    Username:      "your-username",
    Password:      "your-password",
    Sectet2FA:     "your-2fa-secret", // Base32 encoded TOTP secret
    UserAgent:     checkoutOpts.UserAgent,
    AkamaiApiKey:  checkoutOpts.AkamaiApiKey,
    TokenId:       tokenId,
    Proxy:         checkoutOpts.Proxy,
}

// Step 6: Create login client and authenticate
loginClient, err := login.NewLoginClient(loginOpts, &checkoutClient.HttpClient, akamaiAdapter, updateStatus)
if err != nil {
    log.Fatal("Failed to create login client:", err)
}

success, redirectUrl, err := loginClient.Login()
if err != nil {
    log.Fatal("Failed to login:", err)
}

if !success {
    log.Fatal("Login was not successful")
}

// Step 7: Complete checkout
order, err := checkoutClient.Checkout(loginClient, redirectUrl, a2cId, checkoutData)
if err != nil {
    log.Fatal("Failed to complete checkout:", err)
}

log.Println("Checkout completed successfully!")
shr.PPJson(order) // Pretty print order details

```

## Testing

### Run All Tests
```bash
go test ./tests/...
```

### Run Specific Tests
```bash
go test ./tests/ -run TestCart -v
go test ./tests/ -run TestLogin -v  
go test ./tests/ -run TestCheckout -v
```

### Test Runner
```bash
cd tests/runner
go run main.go
```

## Overview

### Key Features

- **Proper TLS Fingerprint** - Mimics real browser behavior, including proper header order and client hello
- **Queue Management** - Handles Best Buy's queue system, from entry to exit into the rest of checkout
- **MFA Support** - TOTP generation for multi-factor auth
- **Start to Finish Checkout Automation** - Capable of checking out any item on Bestbuy with full automation

### Akamai Integration

An Akamai provider is REQUIRED to use this login/checkout api. This is because BestBuy uses Akamai's bot protection, and there are multiple cookies that need to be generated as part of the script. There are multiple API providers out there. One of them is Hyper. Hyper has already been implemented in this project. Feel free to implement other Akamai providers as part of the adapter.

```go
// Production setup with real Akamai API (Hyper SDK)
akamaiAdapter := shr.NewAkamaiAdapterWithHyper(apiKey)
```

## Troubleshooting

### Rare Issues

**Queue Timeouts**
- Add/change queue buffer seconds if needed

**Authentication Failures**
- Check MFA secret format (BestBuy account settings)
- Ensure account credentials are valid

**Akamai Detection**
- Confirm API key is valid
- Check sensor data generation
- Review TLS client configuration
