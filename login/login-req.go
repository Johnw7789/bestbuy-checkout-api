package login

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (l *LoginClient) reqHomePage() (string, error) {
	req, err := http.NewRequest("GET", "https://www.bestbuy.com/", nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {l.Opts.UserAgent},
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
			"priority",
		},
	}

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	abck := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_abck" {
			abck = cookie.Value
			break
		}
	}

	defer resp.Body.Close()

	return abck, nil
}

func (l *LoginClient) reqLoginData() (string, string, error) {
	url := BestbuySigninUrl
	if l.tokenId != "" {
		url = "https://www.bestbuy.com/identity/signin?token=" + l.tokenId
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {l.Opts.UserAgent},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		// "referer":                   {"https://www.bestbuy.com/"},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"accept-language": {"en-US,en;q=0.9"},
		"priority":        {"u=0, i"},
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
			// "referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"priority",
		},
	}

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	l.setBMLSO()

	return string(bodyBytes), resp.Request.URL.String(), nil
}

func (c *LoginClient) setBMLSO() {
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

func (l *LoginClient) reqPublicKey(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":               {"www.bestbuy.com"},
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {l.Opts.UserAgent},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"content-type":       {"application/json"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://www.bestbuy.com/identity/signin?token=" + l.tokenId},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey: {
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"content-type",
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

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqGetTMX(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"User-Agent":         {l.Opts.UserAgent},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"Accept":             {"*/*"},
		"Sec-Fetch-Site":     {"same-site"},
		"Sec-Fetch-Mode":     {"no-cors"},
		"Sec-Fetch-Dest":     {"script"},
		"referer":         {"https://www.bestbuy.com/"},
		"Accept-Encoding": {"gzip, deflate, br, zstd"},
		"Accept-Language": {"en-US,en;q=0.9"},
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
		},
	}

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqGetAkamaiScript(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {l.Opts.UserAgent},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"script"},
		"referer":            {"https://www.bestbuy.com/identity/signin?token=" + l.tokenId},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"accept-language": {"en-US,en;q=0.9"},
		"priority":        {"u=2"},
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

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqSubmitSensorData(url, sensorData string) (string, error) {
	// var data = strings.NewReader(`{"sensor_data":"` + sensorData + `"}`)
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

	// req, err := http.NewRequest("POST", url, data)
	// if err != nil {
	// 	return "", err
	// }

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {l.Opts.UserAgent},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"content-type":       {"text/plain;charset=UTF-8"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://www.bestbuy.com/identity/signin?token=" + l.tokenId},
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

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	abck := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_abck" {
			abck = cookie.Value
			break
		}
	}

	return abck, nil
}

func (l *LoginClient) reqBestbuyLogin(loginJson, encryptedXGrid, xGridB string) (string, error) {
	var data = strings.NewReader(loginJson)

	req, err := http.NewRequest("POST", BestbuyAuthUrl, data)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":               {"www.bestbuy.com"},
		"x-grid-b":           {xGridB},
		"sec-ch-ua-platform": {`"Windows"`},
		"x-grid":             {encryptedXGrid},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {l.Opts.UserAgent},
		"accept":             {"application/json"},
		"content-type":       {"application/json"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://www.bestbuy.com/identity/signin?token=" + l.tokenId},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"x-grid-b",
			"sec-ch-ua-platform",
			"x-grid",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
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

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqBestbuyMFA(mfaJson, encryptedXGrid, xGridB string) (string, error) {
	var data = strings.NewReader(mfaJson)

	req, err := http.NewRequest("POST", BestbuyTwoStepUrl, data)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":               {"www.bestbuy.com"},
		"x-grid-b":           {xGridB},
		"sec-ch-ua-platform": {`"Windows"`},
		"x-grid":             {encryptedXGrid},
		"sec-ch-ua":          {l.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {l.Opts.UserAgent},
		"accept":             {"application/json"},
		"content-type":       {"application/json"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://www.bestbuy.com/identity/signin/twoStepVerification?token=" + l.tokenId},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey: {
			"content-length",
			"x-grid-b",
			"sec-ch-ua-platform",
			"x-grid",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
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

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
