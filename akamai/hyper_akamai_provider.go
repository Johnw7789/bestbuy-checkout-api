package akamai

import (
	"context"

	"github.com/Hyper-Solutions/hyper-sdk-go"
	"github.com/Hyper-Solutions/hyper-sdk-go/akamai"
)

// HyperAkamaiProvider wraps the actual Hyper SDK
type HyperAkamaiProvider struct {
	apiKey  string
	session *hyper.Session
}

// NewHyperAkamaiProvider creates a provider that uses the actual Hyper SDK
func NewHyperAkamaiProvider(apiKey string) *HyperAkamaiProvider {
	return &HyperAkamaiProvider{
		apiKey:  apiKey,
		session: hyper.NewSession(apiKey),
	}
}

// GenerateSensorData calls the actual Hyper SDK
func (h *HyperAkamaiProvider) GenerateSensorData(ctx context.Context, input SensorDataInput) (string, error) {
	hyperInput := &hyper.SensorInput{
		Abck:           input.Abck,
		Bmsz:           input.Bmsz,
		Version:        input.Version,
		UserAgent:      input.UserAgent,
		PageUrl:        input.PageUrl,
		DynamicValues:  input.DynamicValues,
		ScriptHash:     input.ScriptHash,
		IP:             input.IP,
		AcceptLanguage: input.AcceptLanguage,
	}
	
	return h.session.GenerateSensorData(ctx, hyperInput)
}

// GenerateSbsdData calls the actual Hyper SDK
func (h *HyperAkamaiProvider) GenerateSbsdData(ctx context.Context, input SbsdDataInput) (string, error) {
	hyperInput := &hyper.SbsdInput{
		UserAgent:      input.UserAgent,
		Uuid:           input.Uuid,
		PageUrl:        input.PageUrl,
		OCookie:        input.OCookie,
		Script:         input.Script,
		AcceptLanguage: input.AcceptLanguage,
		IP:             input.IP,
	}
	
	return h.session.GenerateSbsdData(ctx, hyperInput)
}

// ParseDynamicValues calls the actual Hyper SDK
func (h *HyperAkamaiProvider) ParseDynamicValues(ctx context.Context, script string) (string, error) {
	hyperInput := &hyper.DynamicInput{Script: script}
	return h.session.ParseV3Dynamic(ctx, hyperInput)
}

// ValidateCookie calls the actual Akamai validation
func (h *HyperAkamaiProvider) ValidateCookie(abck string, requestCount int) bool {
	return akamai.IsCookieValid(abck, requestCount)
}