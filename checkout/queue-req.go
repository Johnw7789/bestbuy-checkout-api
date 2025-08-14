package checkout

import (
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (c *CheckoutClient) reqHomePage() (string, error) {
	req, err := http.NewRequest("GET", "https://www.bestbuy.com/", nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {c.Opts.UserAgent},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey: {
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// * Set the bm_lso cookie, which is almost identical to the bm_so cookie, but with a timestamp appended to it
	c.setBMLSO()

	return string(body), nil
}

func (c *CheckoutClient) reqProdPage() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId), nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {c.Opts.UserAgent},
		"referer":                   {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey: {
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	// if resp.StatusCode == 301 {
	// 	return "", errors.New("Blocked by akamai")
	// }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return string(body), nil
}

func (c *CheckoutClient) requestAddToCart() (*http.Response, error) {
	payload := fmt.Sprintf(`{"items":[{"skuId":"%s"}]}`, c.Opts.SkuId)

	req, err := http.NewRequest(http.MethodPost, BestbuyAddToCartUrl, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		// recaptcha not enabled. static data
		"recaptcha-data":   {`eyJldmVudFV1aWQiOiIxMDMyZWRmYS01MzA3LTQwOWQtODdjNy1hNDU5NmVlYzMzN2IiLCJhY3Rpb24iOiJhZGRUb0NhcnQiLCJlcnJvcnMiOiJSZWNhcHRjaGEgZW5hYmxlZCBjb25maWcgaXMgZmFsc2UuOyBSZWNhcHRjaGEgbm90IGluaXRpYWxpemVkLiBFaXRoZXIgY29uZmlndXJhdGlvbnMgc2V0IHRvIGRpc2FibGVkLCBvciBlcnJvciBwb3B1bGF0aW5nIGNvbmZpZ3VyYXRpb25zLjsgR3JlY2FwdGNoYSBpcyBub3QgZGVmaW5lZC4gQ2Fubm90IGZldGNoIHRva2VuLiJ9`},
		"accept":           {"application/json"},
		"sec-ch-ua":        {c.Opts.UserAgentHint},
		"content-type":     {"application/json; charset=UTF-8"},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {c.Opts.UserAgent},
		"origin":           {"https://www.bestbuy.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":  {"gzip, deflate, br, zstd"},
		"accept-language":  {"en-US,en;q=0.9"},
		"priority":         {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"sec-ch-ua-platform",
			"recaptcha-data",
			"accept",
			"sec-ch-ua",
			"content-type",
			"sec-ch-ua-mobile",
			"user-agent",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	return c.HttpClient.Do(req)
}

func (c *CheckoutClient) requestAddToCartQueue(a2cTransactionCode, a2cTransactionId string) (*http.Response, error) {
	payload := fmt.Sprintf(`{"items":[{"skuId":"%s"}]}`, c.Opts.SkuId)

	req, err := http.NewRequest(http.MethodPost, BestbuyAddToCartUrl, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		// recaptcha not enabled. static data
		"recaptcha-data":            {`eyJldmVudFV1aWQiOiI3ODFmNTVkOS05NTNiLTQzYjktODVmNS02ODg4ODRjOWEwNzkiLCJhY3Rpb24iOiJhZGRUb0NhcnQiLCJlcnJvcnMiOiJSZWNhcHRjaGEgZW5hYmxlZCBjb25maWcgaXMgZmFsc2UuOyBSZWNhcHRjaGEgbm90IGluaXRpYWxpemVkLiBFaXRoZXIgY29uZmlndXJhdGlvbnMgc2V0IHRvIGRpc2FibGVkLCBvciBlcnJvciBwb3B1bGF0aW5nIGNvbmZpZ3VyYXRpb25zLjsgR3JlY2FwdGNoYSBpcyBub3QgZGVmaW5lZC4gQ2Fubm90IGZldGNoIHRva2VuLiJ9`},
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"a2ctransactionreferenceid": {a2cTransactionId},
		"a2ctransactioncode":        {a2cTransactionCode},
		"a2ctransactionwait":        {"undefined"},
		"user-agent":                {c.Opts.UserAgent},
		"accept":                    {"application/json"},
		"content-type":              {"application/json; charset=UTF-8"},
		"origin":                    {"https://www.bestbuy.com"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-dest":            {"empty"},
		"referer":                   {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"sec-ch-ua-platform",
			"recaptcha-data",
			"sec-ch-ua",
			"a2ctransactionreferenceid",
			"a2ctransactioncode",
			"a2ctransactionwait",
			"user-agent",
			"accept",
			"content-type",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	return c.HttpClient.Do(req)
}

func (c *CheckoutClient) requestAddToCartFinal(a2cTransactionId string) (*http.Response, error) {
	payload := fmt.Sprintf(`{"items":[{"skuId":"%s"}]}`, c.Opts.SkuId)

	req, err := http.NewRequest(http.MethodPost, BestbuyAddToCartUrl, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		// recaptcha not enabled. static data
		"recaptcha-data":            {`eyJldmVudFV1aWQiOiJjNTM5NjFhOC1iMWViLTQyNzYtYWVkMy1hMDUyYTc4ZGFhNjgiLCJhY3Rpb24iOiJhZGRUb0NhcnQiLCJlcnJvcnMiOiJSZWNhcHRjaGEgZW5hYmxlZCBjb25maWcgaXMgZmFsc2UuOyBSZWNhcHRjaGEgbm90IGluaXRpYWxpemVkLiBFaXRoZXIgY29uZmlndXJhdGlvbnMgc2V0IHRvIGRpc2FibGVkLCBvciBlcnJvciBwb3B1bGF0aW5nIGNvbmZpZ3VyYXRpb25zLjsgR29vZ2xlIHRva2VuIHJlcXVlc3QgdGltZWQgb3V0LiBUaW1lb3V0PTUwMCJ9`},
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"a2ctransactionreferenceid": {a2cTransactionId},
		"user-agent":                {c.Opts.UserAgent},
		"accept":                    {"application/json"},
		"content-type":              {"application/json; charset=UTF-8"},
		"origin":                    {"https://www.bestbuy.com"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-dest":            {"empty"},
		"referer":                   {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"sec-ch-ua-platform",
			"recaptcha-data",
			"sec-ch-ua",
			"a2ctransactionreferenceid",
			"user-agent",
			"accept",
			"content-type",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	return c.HttpClient.Do(req)
}
