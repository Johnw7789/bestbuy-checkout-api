package login

const (
	BestbuySigninUrl        = "https://www.bestbuy.com/identity/global/signin"
	BestbuyEmailPKUrl       = "https://www.bestbuy.com/api/csiservice/v2/key/cia-email"
	BestbuyCredentialsPKUrl = "https://www.bestbuy.com/api/csiservice/v2/key/credentials-small"
	BestbuyActPKUrl         = "https://www.bestbuy.com/api/csiservice/v2/key/cia-user-activity"
	BestbuyGridPKUrl        = "https://www.bestbuy.com/api/csiservice/v2/key/cia-grid"
	BestbuySubmitEmailUrl   = "https://www.bestbuy.com/identity/password/email"
	BestbuyAuthUrl          = "https://www.bestbuy.com/identity/authenticate"
	BestbuyTwoStepUrl       = "https://www.bestbuy.com/identity/verifyTwoStep"
	BestbuySubmitAuthUrl    = "https://www.bestbuy.com/identity/unlock"
	BestbuyVerifyEmailUrl   = "https://www.bestbuy.com/identity/account/recovery/code"
)

const (
	ErrStepUpRequired = "step up required"
	ErrLoginFailure   = "login failure"
)
