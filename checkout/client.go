package checkout

import (
	"errors"

	tls "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/Johnw7789/bestbuy-checkout/akamai"
)

// * NewCheckoutClient intializes a new http client if one is not provided and returns a new checkout instance
func NewCheckoutClient(opts CheckoutOpts, client *tls.HttpClient, akamaiAdapter *akamai.AkamaiAdapter, updateStatus func(status string)) (*CheckoutClient, error) {
	if opts.AkamaiApiKey == "" {
		return nil, errors.New("AkamaiApiKey is required")
	}

	if opts.UserAgent == "" {
		opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
	}

	if client == nil {
		jar := tls.NewCookieJar()

		options := []tls.HttpClientOption{
			tls.WithTimeoutSeconds(30),
			tls.WithCookieJar(jar),
			tls.WithRandomTLSExtensionOrder(),
			tls.WithClientProfile(profiles.Chrome_133),
		}

		opts.UserAgentHint = `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`

		if opts.Proxy != "" {
			options = append(options, tls.WithProxyUrl(opts.Proxy))
		}

		newClient, err := tls.NewHttpClient(tls.NewNoopLogger(), options...)
		if err != nil {
			return nil, err
		}

		client = &newClient
	}

	return &CheckoutClient{
		HttpClient:    *client,
		UpdateStatus:  updateStatus,
		Opts:          opts,
		AkamaiAdapter: akamaiAdapter,
	}, nil
}
