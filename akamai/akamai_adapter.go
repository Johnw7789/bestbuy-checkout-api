package akamai

import "context"

// AkamaiProvider defines the interface for 3rd party Akamai API providers
type AkamaiProvider interface {
	// GenerateSensorData generates akamai sensor data payload
	GenerateSensorData(ctx context.Context, input SensorDataInput) (string, error)
	
	// GenerateSbsdData generates akamai sbsd data payload  
	GenerateSbsdData(ctx context.Context, input SbsdDataInput) (string, error)
	
	// ParseDynamicValues extracts dynamic values from akamai script
	ParseDynamicValues(ctx context.Context, script string) (string, error)
	
	// ValidateCookie checks if generated cookie is valid
	ValidateCookie(abck string, requestCount int) bool
}

// SensorDataInput contains all data needed to generate sensor data
type SensorDataInput struct {
	Abck           string
	Bmsz           string
	Version        string
	UserAgent      string
	PageUrl        string
	DynamicValues  string
	ScriptHash     string
	IP             string
	AcceptLanguage string
}

// SbsdDataInput contains all data needed to generate sbsd data
type SbsdDataInput struct {
	UserAgent      string
	Uuid           string
	PageUrl        string
	OCookie        string
	Script         string
	AcceptLanguage string
	IP             string
}

// AkamaiAdapter wraps the provider and implements common functionality
type AkamaiAdapter struct {
	provider AkamaiProvider
}

// NewAkamaiAdapter creates a new adapter with the given provider
func NewAkamaiAdapter(provider AkamaiProvider) *AkamaiAdapter {
	return &AkamaiAdapter{
		provider: provider,
	}
}

// NewAkamaiAdapterWithHyper creates a new adapter with the Hyper SDK provider
func NewAkamaiAdapterWithHyper(apiKey string) *AkamaiAdapter {
	provider := NewHyperAkamaiProvider(apiKey)
	return &AkamaiAdapter{
		provider: provider,
	}
}

// GenerateSensorData delegates to the provider
func (a *AkamaiAdapter) GenerateSensorData(ctx context.Context, input SensorDataInput) (string, error) {
	return a.provider.GenerateSensorData(ctx, input)
}

// GenerateSbsdData delegates to the provider
func (a *AkamaiAdapter) GenerateSbsdData(ctx context.Context, input SbsdDataInput) (string, error) {
	return a.provider.GenerateSbsdData(ctx, input)
}

// ParseDynamicValues delegates to the provider
func (a *AkamaiAdapter) ParseDynamicValues(ctx context.Context, script string) (string, error) {
	return a.provider.ParseDynamicValues(ctx, script)
}

// ValidateCookie delegates to the provider
func (a *AkamaiAdapter) ValidateCookie(abck string, requestCount int) bool {
	return a.provider.ValidateCookie(abck, requestCount)
}