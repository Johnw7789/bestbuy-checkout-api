package checkout

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	fhttp "github.com/bogdanfinn/fhttp"
	http "github.com/bogdanfinn/fhttp"
)

func (c *CheckoutClient) reqProdPageRedirect(url string) (string, error) {
	if url == "" {
		url = "https://www.bestbuy.com/site/" + c.Opts.SkuId + ".p?skuId=" + c.Opts.SkuId
	}

	req, err := http.NewRequest("GET", url, nil)
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

	return string(body), nil
}

func (c *CheckoutClient) reqClosestStores() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(BestbuyStoreLocatorUrl, c.Opts.ShippingZipCode), nil)
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

	return string(body), nil
}

// CartCheckout sends a POST request to the cart checkout endpoint
func (c *CheckoutClient) reqCartCheckout(data string) (*http.Response, error) {
	req, err := http.NewRequest("POST", BestbuyGoToCheckoutUrl, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Create headers based on the curl example provided
	req.Header = http.Header{
		"content-length":     {fmt.Sprint(len(data))},
		"user-agent":         {c.Opts.UserAgent},
		"accept":             {"application/json, text/javascript, */*; q=0.01"},
		"accept-language":    {"en-US,en;q=0.9"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"referer":            {BestbuyCartUrl},
		"content-type":       {"application/json"},
		"x-order-id":         {generateUUID()},
		"recaptcha-data":     {"eyJldmVudFV1aWQiOiI4MDMxMzYwYS1kZWIyLTQyZmItYjFiNS1mMWM5YjEzMjNlOGYiLCJhY3Rpb24iOiJnb1RvQ2hlY2tvdXQiLCJlcnJvcnMiOiJSZWNhcHRjaGEgZW5hYmxlZCBjb25maWcgaXMgZmFsc2UuOyBSZWNhcHRjaGEgbm90IGluaXRpYWxpemVkLiBFaXRoZXIgY29uZmlndXJhdGlvbnMgc2V0IHRvIGRpc2FibGVkLCBvciBlcnJvciBwb3B1bGF0aW5nIGNvbmZpZ3VyYXRpb25zLjsgR3JlY2FwdGNoYSBpcyBub3QgZGVmaW5lZC4gQ2Fubm90IGZldGNoIHRva2VuLiJ9"},
		"origin":             {"https://www.bestbuy.com"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"sec-ch-ua":          {c.Opts.UserAgentHint},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"connection":         {"keep-alive"},
		"dnt":                {"1"},
		"sec-gpc":            {"1"},
		"priority":           {"u=0"},
		"te":                 {"trailers"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetOrderDetails fetches the order details from the checkout page
func (c *CheckoutClient) reqOrderDetails() (*http.Response, error) {
	req, err := http.NewRequest("GET", BestbuyCheckoutUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {c.Opts.UserAgent},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"referer":                   {BestbuyCartUrl},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=0, i"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetFulfillment sets the fulfillment method (shipping or pickup)
func (c *CheckoutClient) reqSetFulfillment(data []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf(BestbuyFulfillmentUrl, c.OrderId), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-length":   {fmt.Sprint(len(data))},
		"pragma":           {"no-cache"},
		"cache-control":    {"no-cache"},
		"sec-ch-ua":        {c.Opts.UserAgentHint},
		"accept":           {"application/com.bestbuy.order+json"},
		"x-user-interface": {"DotCom-Optimized"},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {c.Opts.UserAgent},
		"content-type":     {"application/json"},
		"origin":           {"https://www.bestbuy.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {BestbuyShippingEndpoint},
		"accept-encoding":  {"gzip, deflate, br, zstd"},
		"accept-language":  {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetContactInfo sets the contact and shipping information
// func (c *CheckoutClient) reqSetContactInfo(data []byte) (*http.Response, error) {
// 	req, err := http.NewRequest("PATCH", fmt.Sprintf(BestbuyOrderEndpoint, c.OrderId), strings.NewReader(string(data)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header = http.Header{
// 		"content-length":   {fmt.Sprint(len(data))},
// 		"pragma":           {"no-cache"},
// 		"cache-control":    {"no-cache"},
// 		"sec-ch-ua":        {c.Opts.UserAgentHint},
// 		"accept":           {"application/com.bestbuy.order+json"},
// 		"x-user-interface": {"DotCom-Optimized"},
// 		"sec-ch-ua-mobile": {"?0"},
// 		"user-agent":       {c.Opts.UserAgent},
// 		"content-type":     {"application/json"},
// 		"origin":           {"https://www.bestbuy.com"},
// 		"sec-fetch-site":   {"same-origin"},
// 		"sec-fetch-mode":   {"cors"},
// 		"sec-fetch-dest":   {"empty"},
// 		"referer":          {BestbuyShippingEndpoint},
// 		"accept-encoding":  {"gzip, deflate, br, zstd"},
// 		"accept-language":  {"en-US,en;q=0.9"},
// 	}

// 	resp, err := c.HttpClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// RefreshPayment refreshes the payment options
func (c *CheckoutClient) reqRefreshPayment() (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf(BestbuyRefreshPaymentEndpoint, c.OrderId), strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-length":            {"2"},
		"pragma":                    {"no-cache"},
		"cache-control":             {"no-cache"},
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"x-user-interface":          {"DotCom-Optimized"},
		"x-native-checkout-version": {"__VERSION__"},
		"sec-ch-ua-mobile":          {"?0"},
		"user-agent":                {c.Opts.UserAgent},
		"content-type":              {"application/json"},
		"origin":                    {"https://www.bestbuy.com"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-dest":            {"empty"},
		"referer":                   {BestbuyPaymentPageEndpoint},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PaymentPrelookup performs the payment prelookup to get the 3DS reference ID
func (c *CheckoutClient) reqPaymentPrelookup(data []byte, vt string) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf(BestbuyPrelookupEndpoint, c.PaymentId), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = fhttp.Header{
		"content-length":   {fmt.Sprint(len(data))},
		"pragma":           {"no-cache"},
		"cache-control":    {"no-cache"},
		"sec-ch-ua":        {c.Opts.UserAgentHint},
		"vt":               {vt},
		"x-context-id":     {c.OrderId},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {c.Opts.UserAgent},
		"content-type":     {"application/json"},
		"ut":               {"undefined"},
		"x-client":         {"CHECKOUT"},
		"x-request-id":     {fmt.Sprintf("%d", time.Now().UnixNano())},
		"accept":           {"*/*"},
		"origin":           {"https://www.bestbuy.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {BestbuyPaymentPageEndpoint},
		"accept-encoding":  {"gzip, deflate, br, zstd"},
		"accept-language":  {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CheckoutClient) reqPublicKey(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":               {"www.bestbuy.com"},
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {c.Opts.UserAgent},
		"sec-ch-ua":          {c.Opts.UserAgentHint},
		"content-type":       {"application/json"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {BestbuyPaymentPageEndpoint},
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

	resp, err := c.HttpClient.Do(req)
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

// SetPaymentInfo sets the payment information
func (c *CheckoutClient) reqSetPaymentInfo(data []byte) (*http.Response, error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf(BestbuyPaymentEndpoint, c.PaymentId), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-length":   {fmt.Sprint(len(data))},
		"pragma":           {"no-cache"},
		"cache-control":    {"no-cache"},
		"sec-ch-ua":        {c.Opts.UserAgentHint},
		"accept":           {"application/json, text/javascript, */*; q=0.01"},
		"x-client":         {"CHECKOUT"},
		"x-context-id":     {c.OrderId},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {c.Opts.UserAgent},
		"content-type":     {"application/json"},
		"origin":           {"https://www.bestbuy.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {BestbuyPaymentPageEndpoint},
		"accept-encoding":  {"gzip, deflate, br, zstd"},
		"accept-language":  {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CheckoutClient) reqSubmit3DS(jwt string) error {
	formData := fmt.Sprintf("JWT=%s&TermUrl=/payment/r/threeDSecure/redirect&MD=+", jwt)
	req, err := http.NewRequest("POST", "https://centinelapi.cardinalcommerce.com/V2/Cruise/StepUp", strings.NewReader(formData))
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"cache-control":             {"max-age=0"},
		"content-type":              {"application/x-www-form-urlencoded"},
		"origin":                    {"https://www.bestbuy.com"},
		"referer":                   {"https://www.bestbuy.com/"},
		"sec-ch-ua":                 {`"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"iframe"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"cross-site"},
		"sec-fetch-storage-access":  {"none"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {c.Opts.UserAgent},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"priority":                  {"u=1, i"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to submit 3DS")
	}

	return nil
}

func (c *CheckoutClient) reqSubmitCardAuthentication(data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", Bestbuy3DSEndpoint, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-length":   {fmt.Sprint(len(data))},
		"pragma":           {"no-cache"},
		"cache-control":    {"no-cache"},
		"sec-ch-ua":        {c.Opts.UserAgentHint},
		"accept":           {"application/json, text/javascript, */*; q=0.01"},
		"x-client":         {"CHECKOUT"},
		"x-context-id":     {c.OrderId},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {c.Opts.UserAgent},
		"content-type":     {"application/json"},
		"origin":           {"https://www.bestbuy.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {BestbuyPaymentPageEndpoint},
		"accept-encoding":  {"gzip, deflate, br, zstd"},
		"accept-language":  {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
// PlaceOrder places the final order
func (c *CheckoutClient) reqPlaceOrder(data []byte) (*http.Response, error) {
	req, err := fhttp.NewRequest("POST", fmt.Sprintf(BestbuyPlaceOrderEndpoint, c.OrderId), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header = fhttp.Header{
		// "content-length":            {fmt.Sprint(len(data))},
		"pragma":                    {"no-cache"},
		"cache-control":             {"no-cache"},
		"sec-ch-ua":                 {c.Opts.UserAgentHint},
		"accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"x-user-interface":          {"DotCom-Optimized"},
		"x-native-checkout-version": {"__VERSION__"},
		"sec-ch-ua-mobile":          {"?0"},
		"user-agent":                {c.Opts.UserAgent},
		"content-type":              {"application/json"},
		"origin":                    {"https://www.bestbuy.com"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-dest":            {"empty"},
		"referer":                   {BestbuyPaymentPageEndpoint},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
