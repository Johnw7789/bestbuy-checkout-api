package login

import (
	"github.com/Johnw7789/bestbuy-checkout/akamai"
	tls_client "github.com/bogdanfinn/tls-client"
)

type LoginClient struct {
	HttpClient    tls_client.HttpClient
	UpdateStatus  func(status string)
	Opts          LoginOpts
	AkamaiAdapter *akamai.AkamaiAdapter

	tokenId string
}

type LoginOpts struct {
	Username      string
	Password      string
	Sectet2FA     string
	ImapEmail     string
	ImapPassword  string
	Proxy         string
	UserAgent     string
	UserAgentHint string
	AkamaiApiKey  string
	TokenId       string
}

type BestbuyLoginData struct {
	FlowOptions               string
	SocialUserIdFieldName     string
	EncryptedPasswordField    string
	EncryptedAlpha            string
	EmailField                string
	VerificationField         string
	Salmon                    string
	Token                     string
}
