package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Johnw7789/bestbuy-checkout/akamai"
	"github.com/Johnw7789/bestbuy-checkout/checkout"
	"github.com/Johnw7789/bestbuy-checkout/login"

	tls "github.com/bogdanfinn/tls-client"
)

// TestConfig represents the structure of the test configuration file
type TestConfig struct {
	General struct {
		Proxy        string `json:"proxy"`
		UserAgent    string `json:"userAgent"`
		AkamaiApiKey string `json:"akamaiApiKey"`
	} `json:"general"`
	Product struct {
		SkuId string `json:"skuId"`
	} `json:"product"`
	User struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		Secret2FA   string `json:"secret2FA"`
		PhoneNumber string `json:"phoneNumber"`
	} `json:"user"`
	Shipping struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		City       string `json:"city"`
		StateCode  string `json:"stateCode"`
		ZipCode    string `json:"zipCode"`
	} `json:"shipping"`
	Billing struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		City       string `json:"city"`
		StateCode  string `json:"stateCode"`
		ZipCode    string `json:"zipCode"`
	} `json:"billing"`
	Payment struct {
		CardNumber string `json:"cardNumber"`
		ExpMonth   string `json:"expMonth"`
		ExpYear    string `json:"expYear"`
		CVV        string `json:"cvv"`
	} `json:"payment"`
	Testing struct {
		DryRun   bool   `json:"dryRun"`
		TestType string `json:"testType"`
		LogLevel string `json:"logLevel"`
	} `json:"testing"`
}

// LoadTestConfig loads the test configuration from the specified JSON file
func LoadTestConfig(configPath string) (*TestConfig, error) {
	// If configPath is empty, use the default path
	if configPath == "" {
		// Find the config relative to the project root
		configPath = "./runner/test_config.json"
	}

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the JSON configuration
	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// CreateCheckoutClient creates a checkout client using the test configuration
func CreateCheckoutClient(t *testing.T, config *TestConfig) *checkout.CheckoutClient {
	// Set up the checkout options
	opts := checkout.CheckoutOpts{
		SkuId:             config.Product.SkuId,
		AkamaiApiKey:      config.General.AkamaiApiKey,
		UserAgent:         config.General.UserAgent,
		Proxy:             config.General.Proxy,
		Email:             config.User.Email,
		PhoneNumber:       config.User.PhoneNumber,
		ShippingFirstName: config.Shipping.FirstName,
		ShippingLastName:  config.Shipping.LastName,
		ShippingAddress1:  config.Shipping.Address1,
		ShippingAddress2:  config.Shipping.Address2,
		ShippingCity:      config.Shipping.City,
		ShippingStateCode: config.Shipping.StateCode,
		ShippingZipCode:   config.Shipping.ZipCode,
		BillingFirstName:  config.Billing.FirstName,
		BillingLastName:   config.Billing.LastName,
		BillingAddress1:   config.Billing.Address1,
		BillingAddress2:   config.Billing.Address2,
		BillingCity:       config.Billing.City,
		BillingStateCode:  config.Billing.StateCode,
		BillingZipCode:    config.Billing.ZipCode,
		CardNumber:        config.Payment.CardNumber,
		CardExpMonth:      config.Payment.ExpMonth,
		CardExpYear:       config.Payment.ExpYear,
		CVV:               config.Payment.CVV,
	}

	// Create update status function that logs to testing
	updateStatus := func(status string) {
		t.Logf("[Checkout] %s", status)
	}

	// Create Akamai adapter with API key from checkout options
	akamaiAdapter := akamai.NewAkamaiAdapterWithHyper(opts.AkamaiApiKey)

	// Create the checkout client
	client, err := checkout.NewCheckoutClient(opts, nil, akamaiAdapter, updateStatus)
	if err != nil {
		t.Fatalf("Failed to create checkout client: %v", err)
	}

	return client
}

// CreateLoginClient creates a login client using the test configuration
func CreateLoginClient(t *testing.T, config *TestConfig, httpClient *tls.HttpClient, tokenId string) *login.LoginClient {
	// Set up the login options
	opts := login.LoginOpts{
		Username:      config.User.Email,
		Password:      config.User.Password,
		Sectet2FA:     config.User.Secret2FA,
		AkamaiApiKey:  config.General.AkamaiApiKey,
		UserAgent:     config.General.UserAgent,
		TokenId:       tokenId,
		Proxy:         config.General.Proxy,
	}

	// Create update status function that logs to testing
	updateStatus := func(status string) {
		t.Logf("[Login] %s", status)
	}

	// Create Akamai adapter with API key from checkout options
	akamaiAdapter := akamai.NewAkamaiAdapterWithHyper(opts.AkamaiApiKey)

	// Create the login client
	client, err := login.NewLoginClient(opts, httpClient, akamaiAdapter, updateStatus)
	if err != nil {
		t.Fatalf("Failed to create login client: %v", err)
	}

	return client
}
