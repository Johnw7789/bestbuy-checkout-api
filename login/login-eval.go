package login

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/Johnw7789/bestbuy-checkout/encryption"
	"github.com/Johnw7789/bestbuy-checkout/shr"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func (l *LoginClient) getLoginSensorData() (BestbuyLoginData, error) {
	body, currentUrl, err := l.reqLoginData()
	if err != nil {
		return BestbuyLoginData{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return BestbuyLoginData{}, err
	}

	// get script with id = signon-data
	script := doc.Find("script#signon-data")
	if script.Length() == 0 {
		return BestbuyLoginData{}, errors.New("failed to get signon data")
	}

	// Retrieve the script's text content.
	signonData := script.Text()

	var ld BestbuyLoginData
	// * Loop through codeList and alpha list to get encrypted data. we find which one is the correct one by checking if the decoded string contains _X_ or _A_
	passwordArray := gjson.Get(signonData, "codeList")
	for _, passwordField := range passwordArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(passwordField.String())
		if strings.Contains(string(decodedString), "_X_") {
			ld.EncryptedPasswordField = passwordField.String()
			break
		}
	}

	alphaArray := gjson.Get(signonData, "alpha")
	for _, alpha := range alphaArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(shr.Reverse(alpha.String()))
		if strings.Contains(string(decodedString), "_A_") {
			ld.EncryptedAlpha = alpha.String()
			break
		}
	}

	ld.EmailField = gjson.Get(signonData, "emailFieldName").String()
	ld.VerificationField = gjson.Get(signonData, "verificationCodeFieldName").String()
	ld.SocialUserIdFieldName = gjson.Get(signonData, "socialUserIdFieldName").String()
	ld.Salmon = gjson.Get(signonData, "Salmon").String()
	ld.Token = gjson.Get(signonData, "token").String()
	ld.FlowOptions = gjson.Get(signonData, "flowOptions").String()
	zplank := gjson.Get(signonData, "zplankId").String()

	tmxUrl := fmt.Sprintf("https://tmx.bestbuy.com/qla9ftmyheazqr7c.js?%s=ummqowa2&%s=%s", shr.GetRandStr(16), shr.GetRandStr(16), zplank)

	l.tokenId = ld.Token

	_, err = l.reqGetTMX(tmxUrl)
	if err != nil {
		return BestbuyLoginData{}, err
	}

	err = l.handleAkamai(body, currentUrl)
	if err != nil {
		return BestbuyLoginData{}, err
	}

	return ld, nil
}

func (l *LoginClient) getPublicKey(url string) (string, string, error) {
	body, err := l.reqPublicKey(url)
	if err != nil {
		return "", "", err
	}

	publicKey := gjson.Get(body, "publicKey").String()
	keyId := gjson.Get(body, "keyId").String()

	if publicKey == "" || keyId == "" {
		return "", "", errors.New("failed to get public key")
	}

	return publicKey, keyId, nil
}

func (l *LoginClient) getEncData() (encryption.Fingerprint, error) {
	encParams := encryption.EncParams{
		Email:     l.Opts.Username,
		Password:  l.Opts.Password,
		UserAgent: l.Opts.UserAgent,
	}

	var err error

	encParams.EncKeys.EmailKey.PubKey, encParams.EncKeys.EmailKey.KeyId, err = l.getPublicKey(BestbuyEmailPKUrl)
	if err != nil {
		return encryption.Fingerprint{}, err
	}

	encParams.EncKeys.CredentialsKey.PubKey, encParams.EncKeys.CredentialsKey.KeyId, err = l.getPublicKey(BestbuyCredentialsPKUrl)
	if err != nil {
		return encryption.Fingerprint{}, err
	}

	encParams.EncKeys.InfoKey.PubKey, encParams.EncKeys.InfoKey.KeyId, err = l.getPublicKey(BestbuyActPKUrl)
	if err != nil {
		return encryption.Fingerprint{}, err
	}

	encParams.EncKeys.XGridKey.PubKey, encParams.EncKeys.XGridKey.KeyId, err = l.getPublicKey(BestbuyGridPKUrl)
	if err != nil {
		return encryption.Fingerprint{}, err
	}

	fingerprint, err := encryption.GenerateFingerprint(encParams)
	if err != nil {
		return encryption.Fingerprint{}, err
	}

	return fingerprint, nil
}

func (l *LoginClient) getBestbuyLogin(ld BestbuyLoginData, encFingerprint encryption.Fingerprint) (bool, string, error) {
	ld.FlowOptions = "000000010000000000"
	// loginJson := fmt.Sprintf(`{"token":"%s","loginMethod":"CHROME_AUTO","flowOptions":"%s","alpha":"%s","Salmon":"%s","encryptedEmail":"%s","%s":"%s","info":"%s","%s":"%s","recaptchaData": "Error: recaptcha is not enabled."}`, ld.Token, ld.FlowOptions, ld.EncryptedAlpha, ld.Salmon, encFingerprint.EncryptedEmail, ld.EncryptedPasswordField, l.Opts.Password, encFingerprint.EncryptedInfo, ld.EmailField, l.Opts.Username)
	loginJson := fmt.Sprintf(`{"token":"%s","activity":"%s","loginMethod":"UID_PASSWORD","flowOptions":"%s","alpha":"%s","Salmon":"%s","encryptedEmail":"%s","%s":"%s","info":"%s","%s":"%s","recaptchaData": "Error: recaptcha is not enabled."}`, ld.Token, encFingerprint.EncryptedActivity, ld.FlowOptions, ld.EncryptedAlpha, ld.Salmon, encFingerprint.EncryptedEmail, ld.EncryptedPasswordField, l.Opts.Password, encFingerprint.EncryptedInfo, ld.EmailField, l.Opts.Username)

	body, err := l.reqBestbuyLogin(loginJson, encFingerprint.EncryptedXGrid, encFingerprint.XGridB)
	if err != nil {
		return false, "", err
	}

	status := gjson.Get(body, "status").String()
	switch status {
	case "success":
		l.UpdateStatus("Successfully logged in")
		return true, "", nil
	case "stepUpRequired":
		// flowOptions := gjson.Get(body, "flowOptions").String()
		// challengeType := gjson.Get(body, "challengeType").String()

		return false, "", errors.New(ErrStepUpRequired)
	case "failure":
		return false, "", errors.New(ErrLoginFailure)
	default:
		return false, "", errors.New("failed login, status: " + status)
	}
}

func (l *LoginClient) getBestbuyMFA(mfaJson string, encFingerprint encryption.Fingerprint) (bool, string, error) {
	body, err := l.reqBestbuyMFA(mfaJson, encFingerprint.EncryptedXGrid, encFingerprint.XGridB)
	if err != nil {
		return false, "", err
	}

	status := gjson.Get(body, "status").String()
	switch status {
	case "success":
		l.UpdateStatus("Successfully logged in")
		redirectUrl := gjson.Get(body, "redirectUrl").String()
		return true, redirectUrl, nil
	case "stepUpRequired":
		// flowOptions := gjson.Get(body, "flowOptions").String()
		// challengeType := gjson.Get(body, "challengeType").String()

		return false, "", errors.New("step up required")
	case "failure":
		return false, "", errors.New("login failed")
	default:
		return false, "", errors.New("failed login, status: " + status)
	}
}
