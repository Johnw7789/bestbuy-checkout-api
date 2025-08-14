package login

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Johnw7789/bestbuy-checkout/akamai"
	"github.com/Johnw7789/bestbuy-checkout/shr"
)

// * Handles getting valid abck cookie. Posts sensor data up to 3 times in order to get one
func (l *LoginClient) handleAkamai(homePage, currentUrl string) error {
	ctx := context.Background()
	
	if l.AkamaiAdapter == nil {
		return errors.New("akamai adapter not configured")
	}

	ip, err := shr.GetIPAddr()
	if err != nil {
		return errors.New("failed to get IP address")
	}

	scriptPath, err := shr.ParseAkamaiPath(strings.NewReader(homePage))
	if err != nil {
		return err
	}

	akamaiUrl := "https://www.bestbuy.com" + scriptPath

	scriptBody, err := l.reqGetAkamaiScript(akamaiUrl)
	if err != nil {
		return err
	}

	// * Have the akamai api parse the dynamic values from the script
	dynamicVals, err := l.AkamaiAdapter.ParseDynamicValues(ctx, string(scriptBody))
	if err != nil {
		return err
	}

	// * Calc sha256 hash of the script
	scriptHash := sha256.Sum256([]byte(scriptBody))

	// * Submit all 3 sensor data requests to get a valid abck cookie
	return l.submitSensorData(ctx, ip, dynamicVals, scriptHash, akamaiUrl, currentUrl, 3)
}

func (l *LoginClient) submitSensorData(ctx context.Context, ip, vals string, scriptHash [32]byte, akamaiUrl string, currentUrl string, count int) error {
	for reqCount := 1; reqCount <= count; reqCount++ {
		// * Each time we submit sensor data, we need to get the current abck and bmsz cookies again (usually only abck changes, but just in case)
		abck, err := l.getAbck()
		if err != nil {
			continue
		}

		bmsz, err := l.getBmsz()
		if err != nil {
			continue
		}

		l.UpdateStatus("Generating sensor data")
		input := akamai.SensorDataInput{
			Abck:           abck,
			Bmsz:           bmsz,
			Version:        "3",
			UserAgent:      l.Opts.UserAgent,
			PageUrl:        currentUrl,
			DynamicValues:  vals,
			ScriptHash:     fmt.Sprintf("%x", scriptHash),
			IP:             ip,
			AcceptLanguage: "en-US,en;q=0.9",
		}
		
		sensorData, err := l.AkamaiAdapter.GenerateSensorData(ctx, input)
		if err != nil {
			return err
		}

		// l.UpdateStatus("Getting ABCK cookie")

		// * Submit sensor data
		_, err = l.reqSubmitSensorData(akamaiUrl, sensorData)
		if err != nil {
			return err
		}

		// * This no longer works (at least for hyper's api), so we will probably always be sending 3 requests unless they change the validation back to the way it was
		valid := l.AkamaiAdapter.ValidateCookie(abck, reqCount)
		if valid {
			break
		}
	}

	return nil
}

func (l *LoginClient) getAbck() (string, error) {
	abck := ""
	u, _ := url.Parse("https://www.bestbuy.com")

	for _, cook := range l.HttpClient.GetCookieJar().Cookies(u) {
		if cook.Name == "_abck" {
			abck = cook.Value
			return abck, nil
		}
	}

	return abck, nil
}

func (l *LoginClient) getBmsz() (string, error) {
	bmsz := ""
	u, _ := url.Parse("https://www.bestbuy.com")
	for _, cook := range l.HttpClient.GetCookieJar().Cookies(u) {
		if cook.Name == "bm_sz" {
			bmsz = cook.Value
			return bmsz, nil
		}
	}

	return bmsz, nil
}
