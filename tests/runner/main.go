package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
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

func main() {
	// Define command line flags
	configPath := flag.String("config", "", "Path to test config file (default: config/test_config.json)")
	testType := flag.String("test", "all", "Test to run: login, cart, checkout, dry-run, all")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override test type if specified via command line
	if *testType != "all" {
		config.Testing.TestType = *testType
	}

	fmt.Printf("Running %s tests...\n", config.Testing.TestType)

	// Build test command based on configuration
	var testArgs []string
	var exitCode int

	// Add verbose flag if specified
	if *verbose {
		testArgs = append(testArgs, "-v")
	}

	// Run appropriate tests based on configuration
	switch config.Testing.TestType {
	case "login":
		fmt.Println("=== Running Login Test ===")
		testArgs = append(testArgs, "-run", "TestLogin")
		exitCode = runGoTest(testArgs)

	case "cart":
		fmt.Println("=== Running Cart Test ===")
		testArgs = append(testArgs, "-run", "TestCartItem")
		exitCode = runGoTest(testArgs)

	case "checkout":
		fmt.Println("=== Running Full Checkout Test ===")
		// Set DryRun to false to actually place an order
		config.Testing.DryRun = false
		saveConfig(*configPath, config)
		
		testArgs = append(testArgs, "-run", "TestCheckoutFlow")
		exitCode = runGoTest(testArgs)

	case "dry-run":
		fmt.Println("=== Running Checkout Dry Run Test ===")
		// Ensure DryRun is true to prevent placing a real order
		config.Testing.DryRun = true
		saveConfig(*configPath, config)
		
		testArgs = append(testArgs, "-run", "TestCheckoutDryRun")
		exitCode = runGoTest(testArgs)

	default:
		fmt.Printf("Unknown test type: %s\n", config.Testing.TestType)
		fmt.Println("Available test types: login, cart, checkout, dry-run, all")
		exitCode = 1
	}

	os.Exit(exitCode)
}

// runGoTest executes the go test command with the provided arguments
func runGoTest(args []string) int {
	// Build the command to run the tests
	baseArgs := []string{"test", "../..."}
	allArgs := append(baseArgs, args...)
	
	cmd := exec.Command("go", allArgs...)
	
	// Connect command output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Run the command
	err := cmd.Run()
	
	// Check for command execution errors
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		fmt.Printf("Error running tests: %v\n", err)
		return 1
	}
	
	return 0
}

func loadConfig(configPath string) (*TestConfig, error) {
	// If path is not specified, use default in the current directory
	if configPath == "" {
		configPath = "test_config.json"
	}

	fmt.Printf("Loading config from: %s\n", configPath)

	// Try to read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If not found, try in the project root directory
		fmt.Printf("Could not read config from %s, trying in project root...\n", configPath)
		rootPath := "../../config/test_config.json"
		data, err = os.ReadFile(rootPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file from %s or %s: %w", configPath, rootPath, err)
		}
		fmt.Printf("Successfully loaded config from %s\n", rootPath)
	} else {
		fmt.Printf("Successfully loaded config from %s\n", configPath)
	}

	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func saveConfig(configPath string, config *TestConfig) error {
	// If path is not specified, use default
	if configPath == "" {
		configPath = "test_config.json"
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
