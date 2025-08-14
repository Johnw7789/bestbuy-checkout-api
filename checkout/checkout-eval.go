package checkout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func (c *CheckoutClient) getOrderDetails() (bool, error) {
	resp, err := c.reqOrderDetails()
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("failed to get checkout page, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Extract the order data using regex
	re := regexp.MustCompile(`var orderData = (.*?);`)
	matches := re.FindSubmatch([]byte(body))
	if len(matches) < 2 {
		return false, fmt.Errorf("could not find order data in checkout page")
	}

	orderData := matches[1]
	c.OrderId = gjson.Get(string(orderData), "id").String()
	c.PaymentId = gjson.Get(string(orderData), "payment.id").String()
	c.UserOrderId = gjson.Get(string(orderData), "userOrderId").String()
	shippingAvailable := gjson.Get(string(orderData), "items.0.meta.fulfillmentEligibilities.shippingEligible").Bool()
	// c.ItemIds = make([]string, len(orderData.Items))
	// for i, item := range orderData.Items {
	// 	c.ItemIds[i] = item.ID
	// }

	c.UpdateStatus(fmt.Sprintf("Order details retrieved: OrderID=%s, PaymentID=%s", c.OrderId, c.PaymentId))
	return shippingAvailable, nil
}

func (c *CheckoutClient) getSetFufillment() error {
	var fulfillmentPayload FulfillmentPayload
	for _, itemID := range c.ItemIds {
		// if c.Opts.LocationID == "" {
		// Use shipping
		fulfillmentPayload.Items = append(fulfillmentPayload.Items, ItemFulfillment{
			Id:   itemID,
			Type: "default",
			Fulfillment: Fulfillment{
				Address: ShipAddress{
					Addressline1:         c.Opts.ShippingAddress1,
					Addressline2:         c.Opts.ShippingAddress2,
					City:                 c.Opts.ShippingCity,
					Classifer:            "RESIDENTIAL",
					Country:              "US",
					Firstname:            c.Opts.ShippingFirstName,
					Lastname:             c.Opts.ShippingLastName,
					Middleinitial:        "",
					Phonenumber:          c.Opts.PhoneNumber,
					Postalcode:           c.Opts.ShippingZipCode,
					SaveAddressAsDefault: false,
					Savedtoprofile:       false,
					State:                c.Opts.ShippingStateCode,
					Type:                 "shipping",
					Useforbilling:        false,
				},
				Type: "shipping",
			},
		})
		// } else {
		// 	// Use in-store pickup

		// }
	}

	data, err := json.Marshal(fulfillmentPayload)
	if err != nil {
		return err
	}

	resp, err := c.reqSetFulfillment(data)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to set fulfillment method, status: %d", resp.StatusCode)
	}

	c.UpdateStatus("Fulfillment method set successfully")
	return nil
}

func (c *CheckoutClient) getSetContactInfo() error {
	if len(c.ItemIds) == 0 {
		return errors.New("no items found in order")
	}

	// Create the request payload
	// contactInfo := ContactInfo{
	// 	Phonenumber:     c.Opts.PhoneNumber,
	// 	Smsnotifynumber: "",
	// 	Smsoptin:        false,
	// 	Emailaddress:    c.Opts.Email,
	// 	// Items:           []ShipItem{},
	// }

	// for _, itemID := range c.ItemIds {
	// 	if c.Opts.LocationID == "" {
	// 		// Use shipping
	// 		contactInfo.Items = append(contactInfo.Items, ShipItem{
	// 			ID:   itemID,
	// 			Type: "DEFAULT",
	// 			Selectedfulfillment: Selectedfulfillment{
	// 				Shipping: Shipping{
	// 					Address: ShipAddress{
	// 						Country:             c.Opts.ShippingCountryCode,
	// 						Street2:             c.Opts.ShippingAddress2,
	// 						Useaddressasbilling: false,
	// 						Lastname:            c.Opts.ShippingLastName,
	// 						Street:              c.Opts.ShippingAddress1,
	// 						City:                c.Opts.ShippingCity,
	// 						Zipcode:             c.Opts.ShippingZipCode,
	// 						State:               c.Opts.ShippingStateCode,
	// 						Firstname:           c.Opts.ShippingFirstName,
	// 						Dayphonenumber:      c.Opts.PhoneNumber,
	// 						Type:                "RESIDENTIAL",
	// 					},
	// 				},
	// 			},
	// 			Giftmessageselected: false,
	// 		})
	// 	} else {
	// 		// For pickup, we still need to include the item but don't need shipping address
	// 		contactInfo.Items = append(contactInfo.Items, ShipItem{
	// 			ID:                  itemID,
	// 			Type:                "DEFAULT",
	// 			Giftmessageselected: false,
	// 		})
	// 	}
	// }

	// data, err := json.Marshal(contactInfo)
	// if err != nil {
	// 	return err
	// }

	// resp, err := c.reqSetContactInfo(data)
	// if err != nil {
	// 	return err
	// }

	// if resp.StatusCode != 200 {
	// 	return fmt.Errorf("failed to set contact information, status: %d", resp.StatusCode)
	// }

	c.UpdateStatus("Contact and shipping information set successfully")
	return nil
}

func (c *CheckoutClient) getPublicKey(url string) (string, string, error) {
	body, err := c.reqPublicKey(url)
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

func (c *CheckoutClient) getPaymentPrelookup() error {
	payload := PrelookupRequest{
		BinNumber: c.Opts.CardNumber[0:6],
		Browserinfo: Browserinfo{
			Javaenabled: false,
			Language:    "en-US",
			Useragent:   c.Opts.UserAgent,
			Height:      "1440",
			Width:       "2560",
			Timezone:    "240",
			Colordepth:  "24",
		},
		Orderid:             c.UserOrderId,
		PaymentId:           c.PaymentId,
		PayComponentVersion: "STANDARD_REDESIGN",
		PayVersion:          "5.3.276",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	vt := ""
	u, _ := url.Parse("https://www.bestbuy.com")
	for _, cookie := range c.HttpClient.GetCookieJar().Cookies(u) {
		if cookie.Name == "vt" {
			vt = cookie.Value
			break
		}
	}

	resp, err := c.reqPaymentPrelookup(data, vt)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to prelookup payment, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 	// Parse the response to get the 3DS reference ID
	var prelookupResp PrelookupResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&prelookupResp)
	if err != nil {
		return err
	}

	c.ThreeDsId = prelookupResp.ThreeDSReferenceId
	c.PaymentId = prelookupResp.PaymentId

	return nil
}

func (c *CheckoutClient) getSetPaymentInfo(num string) (string, error) {
	if c.PaymentId == "" {
		return "", fmt.Errorf("payment ID not found")
	}

	// Get card type based on first digit
	cardType := getCardType(c.Opts.CardNumber)

	billing := Billingaddress{
		Addressline1:    c.Opts.BillingAddress1,
		Addressline2:    c.Opts.BillingAddress2,
		City:            c.Opts.BillingCity,
		Country:         "US",
		Dayphone:        c.Opts.PhoneNumber,
		Emailaddress:    c.Opts.Email,
		Firstname:       c.Opts.BillingFirstName,
		IsSimpleAddress: false,
		Lastname:        c.Opts.BillingLastName,
		Postalcode:      c.Opts.BillingZipCode,
		Standardized:    true,
		State:           c.Opts.BillingStateCode,
		Useroverridden:  false,
	}

	paymentReq := PaymentRequest{
		Billingaddress: billing,
		Creditcard: Creditcard{
			AppliedPoints:       "",
			Binnumber:           c.Opts.CardNumber[0:6],
			CreditCardProfileId: "",
			Cvv:                 c.Opts.CVV,
			Default:             false,
			Expmonth:            c.Opts.CardExpMonth,
			Expyear:             c.Opts.CardExpYear,
			Number:              num,
			Orderid:             c.UserOrderId,
			PurchaseOrderNumber: "",
			Savetoprofile:       false,
			Type:                cardType,
			Virtualcard:         false,
		},
	}

	data, err := json.Marshal(paymentReq)
	if err != nil {
		return "", err
	}

	resp, err := c.reqSetPaymentInfo(data)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to set payment information, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return gjson.Get(string(body), "creditCard.number").String(), nil
}

func (c *CheckoutClient) getRefreshPayment() error {
	resp, err := c.reqRefreshPayment()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to refresh payment options, status: %d", resp.StatusCode)
	}

	return nil
}

func (c *CheckoutClient) getSubmitCardAuthentication() error {
	cardAuthPayload := CardAuthentication{
		TransactionId: uuid.New().String(), // 3DS flow just generates a random uuid for some reason
		OrderId: 	 c.OrderId,
	}

	data, err := json.Marshal(cardAuthPayload)
	if err != nil {
		return err 
	}

	resp, err := c.reqSubmitCardAuthentication(data)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to submit card authentication, status: %d", resp.StatusCode)
	}

	return nil
}

func (c *CheckoutClient) getPlaceOrder() (Order, error) {
	placeOrderReq := PlaceOrderRequest{
		Metadata: Metadata{
			Browserinfo: Browserinfo{
				Javaenabled: false,
				Language:    "en-US",
				Useragent:   c.Opts.UserAgent,
				Height:      "1440",
				Width:       "2560",
				Timezone:    "240",
				Colordepth:  "24",
			},
			Buylegalmessaging: []string{},
			Userdevicedetails: Userdevicedetails{
				ForterToken: fmt.Sprintf("93813b86f7ce44fd9838cc96f03de6fe_%d__UDF43-m4_19ck__tt", time.Now().UnixMilli()),
			},
		},
		User: User{},
	}

	data, err := json.Marshal(placeOrderReq)
	if err != nil {
		return Order{}, err
	}

	resp, err := c.reqPlaceOrder(data)
	if err != nil {
		return Order{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Order{}, err
	}

	// * Handle 3DS
	if resp.StatusCode == 412 && strings.Contains(string(body), "PAY_SECURE_REDIRECT") {
		c.UpdateStatus("Handling 3DS")

		threeDSJWT := gjson.GetBytes(body, "errors.0.paySecureResponse.stepUpJwt").String()

		err = c.reqSubmit3DS(threeDSJWT)
		if err != nil {
			return Order{}, err
		}

		// After 3DS, we need to sleep for the 3DS redirect interval (90 seconds) in order to bypass bank verification
		time.Sleep(90 * time.Second)

		err = c.getSubmitCardAuthentication()
		if err != nil {
			return Order{}, err
		}

		c.UpdateStatus("Submitting payment")
		resp, err = c.reqPlaceOrder(data)
		if err != nil {
			return Order{}, err
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return Order{}, err
		}
	}

	if resp.StatusCode != 200 {
		return Order{}, fmt.Errorf("failed to place order, status: %d", resp.StatusCode)
	}

	status := gjson.GetBytes(body, "status").String()
	if !strings.EqualFold(status, "submitted") {
		return Order{}, fmt.Errorf("failed to place order, status: %s", status)
	}

	return Order{
		OrderId:    gjson.GetBytes(body, "userOrderId").String(),
		TotalPrice: gjson.GetBytes(body, "price.orderTotal").Float(),
		PickupDate: gjson.GetBytes(body, "items.0.fulfillment.pickupDate").String(),
	}, nil
}
