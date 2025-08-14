package tests

import (
	"testing"
)

// TestLogin tests the login functionality
func TestLogin(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping login test in short mode")
	}

	// Load test configuration
	config, err := LoadTestConfig("")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	t.Log("Starting login test")

	// Create login client (with empty token ID for direct login)
	loginClient := CreateLoginClient(t, config, nil, "")

	// Attempt to login
	success, redirectURL, err := loginClient.Login()

	// Check for errors
	if err != nil {
		t.Fatalf("Login failed with error: %v", err)
	}

	// Check login success
	if !success {
		t.Error("Login was not successful")
	}

	// Verify we got a redirect URL
	if redirectURL == "" {
		t.Error("Expected redirect URL but got empty string")
	}

	t.Logf("Login successful, redirect URL: %s", redirectURL)
}
