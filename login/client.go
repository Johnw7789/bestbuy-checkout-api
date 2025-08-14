package login

import (
	"errors"

	"github.com/Johnw7789/bestbuy-checkout/akamai"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// * NewLoginClient intializes a new http client if one is not provided and returns a new login instance
func NewLoginClient(opts LoginOpts, client *tls_client.HttpClient, akamaiAdapter *akamai.AkamaiAdapter, updateStatus func(status string)) (*LoginClient, error) {
	if opts.AkamaiApiKey == "" {
		return nil, errors.New("AkamaiApiKey is required")
	}

	if opts.UserAgent == "" {
		opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	}

	if client == nil {
		jar := tls_client.NewCookieJar()

		options := []tls_client.HttpClientOption{
			tls_client.WithTimeoutSeconds(30),
			tls_client.WithCookieJar(jar),
			tls_client.WithRandomTLSExtensionOrder(),
			tls_client.WithClientProfile(profiles.Chrome_133),
		}

		opts.UserAgentHint = `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`

		if opts.Proxy != "" {
			options = append(options, tls_client.WithProxyUrl(opts.Proxy))
		}

		newClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
		if err != nil {
			return nil, err
		}

		client = &newClient
	}

	return &LoginClient{
		HttpClient:    *client,
		UpdateStatus:  updateStatus,
		Opts:          opts,
		AkamaiAdapter: akamaiAdapter,
		tokenId:       opts.TokenId,
	}, nil
}
