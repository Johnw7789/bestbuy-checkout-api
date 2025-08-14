package tests

import (
	"testing"
)

// TestCartItem tests adding an item to the cart
func TestCartItem(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping cart test in short mode")
	}

	// Load test configuration
	config, err := LoadTestConfig("")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	t.Log("Starting cart test")

	// Create checkout client
	checkoutClient := CreateCheckoutClient(t, config)

	// Add item to cart
	tokenId, a2cId, checkoutData, err := checkoutClient.CartItem()

	// We expect either success or a "no queue found" error (which is not a failure)
	if err != nil && err.Error() != "no queue found" {
		t.Fatalf("Cart process failed: %v", err)
	}

	// If we had no queue but got checkout data, that's success
	if err != nil && err.Error() == "no queue found" {
		if checkoutData == "" {
			t.Error("Got 'no queue found' but checkout data is also empty")
		} else {
			t.Logf("No queue found, but item was added to cart successfully")
			t.Logf("Checkout data: %s", checkoutData)
		}
		return
	}

	if tokenId == "" {
		t.Error("Expected tokenId but got empty string")
	}
	
	if a2cId == "" {
		t.Error("Expected a2cId but got empty string")
	}
	
	if checkoutData == "" {
		t.Error("Expected checkoutData but got empty string")
	}

	t.Logf("TokenId: %s", tokenId)
	t.Logf("A2CId: %s", a2cId)
	t.Logf("Cart process completed successfully")
}
