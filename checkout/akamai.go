package checkout

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Johnw7789/bestbuy-checkout/akamai"
	"github.com/Johnw7789/bestbuy-checkout/shr"
	http "github.com/bogdanfinn/fhttp"
)

// * Handles getting valid abck cookie. Posts sensor data up to 3 times in order to get one
func (c *CheckoutClient) handleAkamai(homePage string, doSbsd bool) error {
	ctx := context.Background()
	
	if c.AkamaiAdapter == nil {
		return errors.New("akamai adapter not configured")
	}

	ip, err := shr.GetIPAddr()
	if err != nil {
		return errors.New("failed to get IP address")
	}

	if doSbsd {
		// * Since the path is constantly changing every few minutes, we have to parse it from the body
		sbsdFullPath, sbsdPath, uuid, err := shr.ParseSBSD(strings.NewReader(homePage))
		if err != nil {
			return err
		}

		sbsdUrl := "https://www.bestbuy.com" + sbsdFullPath

		// c.UpdateStatus("Getting SBSD script")
		sbsdScriptBody, err := c.reqGetAkamaiScript(sbsdUrl)
		if err != nil {
			return err
		}

		// // * Regardless of if we get sbsd script blocking the page, submit the sbsd data anyway
		err = c.submitSbsdData(ctx, ip, string(sbsdScriptBody), uuid, "https://www.bestbuy.com"+sbsdPath, 3)
		if err != nil {
			return err
		}
	}

	c.UpdateStatus("Getting product page")

	// * Now we need to get the akamai script from the prod page
	prodPage, err := c.reqProdPage()
	if err != nil {
		return err
	}

	scriptPath, err := shr.ParseAkamaiPath(strings.NewReader(prodPage))
	if err != nil {
		return err
	}

	akamaiUrl := "https://www.bestbuy.com" + scriptPath

	c.UpdateStatus("Getting Akamai script")
	scriptBody, err := c.reqGetAkamaiScript(akamaiUrl)
	if err != nil {
		return err
	}

	// * Have the akamai api parse the dynamic values from the script
	dynamicVals, err := c.AkamaiAdapter.ParseDynamicValues(ctx, string(scriptBody))
	if err != nil {
		return err
	}

	// * Calc sha256 hash of the script
	scriptHash := sha256.Sum256([]byte(scriptBody))

	// * Submit all 3 sensor data requests to get a valid abck cookie
	return c.submitSensorData(ctx, ip, dynamicVals, scriptHash, akamaiUrl, 3)
}

func (c *CheckoutClient) submitSbsdData(ctx context.Context, ip, script, uuid, sbsdUrl string, count int) error {
	oCook, err := c.getOCookie()
	if err != nil {
		return err
	}

	for reqCount := 1; reqCount <= count; reqCount++ {
		c.UpdateStatus("Generating SBSD data")
		input := akamai.SbsdDataInput{
			UserAgent:      c.Opts.UserAgent,
			Uuid:           uuid,
			PageUrl:        fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId),
			OCookie:        oCook,
			Script:         script,
			AcceptLanguage: "en-US,en;q=0.9",
			IP:             ip,
		}
		
		payload, err := c.AkamaiAdapter.GenerateSbsdData(ctx, input)
		if err != nil {
			return err
		}

		_, err = c.reqSubmitSBSD(sbsdUrl, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CheckoutClient) submitSensorData(ctx context.Context, ip, vals string, scriptHash [32]byte, akamaiUrl string, count int) error {
	for reqCount := 1; reqCount <= count; reqCount++ {
		// * Each time we submit sensor data, we need to get the current abck and bmsz cookies again (usually only abck changes, but just in case)
		abck, err := c.getAbck()
		if err != nil {
			continue
		}

		bmsz, err := c.getBmsz()
		if err != nil {
			continue
		}

		c.UpdateStatus("Generating sensor data")
		input := akamai.SensorDataInput{
			Abck:           abck,
			Bmsz:           bmsz,
			Version:        "3",
			UserAgent:      c.Opts.UserAgent,
			PageUrl:        fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId),
			DynamicValues:  vals,
			ScriptHash:     fmt.Sprintf("%x", scriptHash),
			IP:             ip,
			AcceptLanguage: "en-US,en;q=0.9",
		}
		
		sensorData, err := c.AkamaiAdapter.GenerateSensorData(ctx, input)
		if err != nil {
			return err
		}

		// c.UpdateStatus("Getting ABCK cookie")

		// * Submit sensor data
		_, err = c.reqSubmitSensorData(akamaiUrl, sensorData)
		if err != nil {
			return err
		}

		// * This no longer works (at least for hyper's api), so we will probably always be sending 3 requests unless they change the validation back to the way it was
		valid := c.AkamaiAdapter.ValidateCookie(abck, reqCount)
		if valid {
			break
		}
	}

	return nil
}

func (c *CheckoutClient) setBMLSO() {
	u, _ := url.Parse("https://www.bestbuy.com")
	for _, cook := range c.HttpClient.GetCookieJar().Cookies(u) {
		if cook.Name == "bm_so" {
			// * Set the bm_lso cookie to the same value as bm_so, but with the current time in milliseconds appended to it
			ms := time.Now().UnixMilli() // e.g. 1746574295059

			msStr := strconv.FormatInt(ms, 10)

			c.HttpClient.GetCookieJar().SetCookies(u, []*http.Cookie{
				{
					Name:  "bm_lso",
					Value: cook.Value + "^" + msStr,
				},
			})

			break
		}
	}
}

// * Get the bm_so cookie from the cookie jar needed for the akamai sbsd payload
func (c *CheckoutClient) getOCookie() (string, error) {
	oCook := ""
	u, _ := url.Parse("https://www.bestbuy.com")

	for _, cook := range c.HttpClient.GetCookieJar().Cookies(u) {
		if cook.Name == "bm_so" {
			oCook = cook.Value
			return oCook, nil
		}
	}

	return oCook, nil
}

func (c *CheckoutClient) getAbck() (string, error) {
	abck := ""
	u, _ := url.Parse("https://www.bestbuy.com")

	for _, cook := range c.HttpClient.GetCookieJar().Cookies(u) {
		if cook.Name == "_abck" {
			abck = cook.Value
			return abck, nil
		}
	}

	return abck, nil
}

func (l *CheckoutClient) getBmsz() (string, error) {
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
