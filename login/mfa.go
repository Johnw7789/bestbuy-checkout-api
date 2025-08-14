package login

import (
	"errors"
	"fmt"

	"github.com/Johnw7789/bestbuy-checkout/encryption"
)

// * Handles the TOTP 2FA flow
func (l *LoginClient) handleMFA(ld BestbuyLoginData, encFingerprint encryption.Fingerprint) (bool, string, error) {
	if l.Opts.Sectet2FA == "" {
		return false, "", errors.New("2FA secret is required")
	}

	otp, err := getOTP(l.Opts.Sectet2FA)
	if err != nil {
		return false, "", errors.New("failed to generate OTP")
	}

	l.UpdateStatus("Generated OTP: " + otp)

	// * Flow options for: don't trust device
	flowOptions := "00500001210000001"
	mfaJson := fmt.Sprintf(`{"token":"%s","flowOptions":"%s","password":"%s","socialUserIdFieldName":"%s","mt":null,"%s":"%s","%s":"%s"}`, ld.Token, flowOptions, encFingerprint.EncryptedPassword, ld.SocialUserIdFieldName, ld.EmailField, l.Opts.Username, ld.VerificationField, otp)

	l.UpdateStatus("Submitting OTP")
	return l.getBestbuyMFA(mfaJson, encFingerprint)
}
