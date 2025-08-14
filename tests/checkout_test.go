package tests

import (
	"testing"
)

// TestCheckoutFlow tests the complete checkout flow
func TestCheckoutFlow(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping checkout test in short mode")
	}

	// Load test configuration
	config, err := LoadTestConfig("")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Check if we're running in dry run mode
	if config.Testing.DryRun {
		t.Skip("Skipping checkout completion in dry run mode to avoid placing a real order. Set dryRun: false in config to enable.")
	}

	t.Log("Starting checkout flow test")

	// Step 1: Create checkout client
	checkoutClient := CreateCheckoutClient(t, config)

	// Step 2: Add to cart and get necessary tokens
	tokenId, a2cId, checkoutData, err := checkoutClient.CartItem()
	if err != nil && err.Error() != "no queue found" {
		t.Fatalf("Failed to add to cart: %v", err)
	}

	// If we had no queue, we should still have checkout data
	if err != nil && err.Error() == "no queue found" {
		if checkoutData == "" {
			t.Fatalf("Got 'no queue found' but checkout data is also empty")
		}
		t.Logf("No queue found, but item was added to cart successfully")
	} else {
		t.Logf("Queue was processed successfully")
		t.Logf("TokenId: %s", tokenId)
		t.Logf("A2CId: %s", a2cId)
	}

	// Step 3: Create login client and login
	loginClient := CreateLoginClient(t, config, &checkoutClient.HttpClient, tokenId)
	
	// Attempt to login
	success, redirectUrl, err := loginClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if !success {
		t.Fatalf("Login was not successful")
	}

	t.Logf("Login successful, redirect URL: %s", redirectUrl)

	// Step 4: Complete checkout process
	t.Log("Starting checkout process")
	order, err := checkoutClient.Checkout(loginClient, redirectUrl, a2cId, checkoutData)
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	t.Logf("Order placed successfully: OrderID=%s, Total=$%.2f", order.OrderId, order.TotalPrice)
}

// TestCheckoutDryRun tests the checkout flow but stops short of placing the order
func TestCheckoutDryRun(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping checkout dry run test in short mode")
	}

	// Load test configuration
	config, err := LoadTestConfig("")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	t.Log("Starting checkout dry run test")

	// Step 1: Create checkout client
	checkoutClient := CreateCheckoutClient(t, config)

	// Step 2: Add to cart and get necessary tokens
	tokenId, a2cId, checkoutData, err := checkoutClient.CartItem()
	if err != nil && err.Error() != "no queue found" {
		t.Fatalf("Failed to add to cart: %v", err)
	}

	// If we had no queue, we should still have checkout data
	if err != nil && err.Error() == "no queue found" {
		if checkoutData == "" {
			t.Fatalf("Got 'no queue found' but checkout data is also empty")
		}
		t.Logf("No queue found, but item was added to cart successfully")
	} else {
		t.Logf("Queue was processed successfully")
		t.Logf("TokenId: %s", tokenId)
		t.Logf("A2CId: %s", a2cId)
	}

	// Step 3: Create login client and login
	loginClient := CreateLoginClient(t, config, &checkoutClient.HttpClient, tokenId)
	
	// Attempt to login
	success, redirectUrl, err := loginClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if !success {
		t.Fatalf("Login was not successful")
	}

	t.Logf("Login successful, redirect URL: %s", redirectUrl)

	// For a dry run, we stop here instead of completing the checkout
	t.Log("Checkout dry run completed successfully. Order was not placed.")
}
