package checkout

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (c *CheckoutClient) reqGetAkamaiScript(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {c.Opts.UserAgent},
		"sec-ch-ua":          {c.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"script"},
		"referer":            {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"u=2"},
		http.HeaderOrderKey: {
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"accept",
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

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyScript, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyScript), nil
}

func (c *CheckoutClient) reqSubmitSensorData(url, sensorData string) (string, error) {
	type SensorData struct {
		SensorData string `json:"sensor_data"`
	}

	data, err := json.Marshal(SensorData{SensorData: sensorData})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {c.Opts.UserAgent},
		"sec-ch-ua":          {c.Opts.UserAgentHint},
		"content-type":       {"text/plain;charset=UTF-8"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US;en;q=0.9"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"content-type",
			"sec-ch-ua-mobile",
			"accept",
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

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return "", nil
}

func (c *CheckoutClient) reqSubmitSBSD(url, body string) (*http.Response, error) {
	type Body struct {
		Body string `json:"body"`
	}

	data, err := json.Marshal(Body{Body: body})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {c.Opts.UserAgent},
		"sec-ch-ua":          {c.Opts.UserAgentHint},
		"content-type":       {"application/json"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {fmt.Sprintf("https://www.bestbuy.com/site/%s.p", c.Opts.SkuId)},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"content-type",
			"sec-ch-ua-mobile",
			"accept",
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
