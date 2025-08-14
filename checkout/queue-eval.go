package checkout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

// * Here are the cookies absolutely required to pass queue (in browser):
//  - _abck
//  - bm_s
//  - bm_sz
//  - CTT
//  - SID
//  - vt

// * Add to cart and return the a2cTransactionCode and a2cTransactionId response headers
func (c *CheckoutClient) tryCart() (string, string, string, error) {
	c.UpdateStatus("Adding to cart")

	resp, err := c.requestAddToCart()
	if err != nil {
		return "", "", "", err
	}

	// * Parse the response to get the a2cTransactionCode and a2cTransactionCode
	a2cTransactionCode := resp.Header.Get("a2ctransactioncode")
	a2cTransactionId := resp.Header.Get("a2ctransactionreferenceid")

	// * If the a2cTransactionCode and a2cTransactionId are empty, it means there is no queue
	if a2cTransactionCode == "" && a2cTransactionId == "" && (resp.StatusCode == 200 || resp.StatusCode == 201) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", "", err
		}

		defer resp.Body.Close()

		productTotalStr := fmt.Sprintf("%.2f", gjson.Get(string(body), "cartSubTotal").Float())
		productTotalFormatted := strings.Replace(productTotalStr, ".", ",", 1)
		cartCount := gjson.Get(string(body), "cartCount").Int()
		cartSubTotal := gjson.Get(string(body), "cartSubTotal").Float()

		c.UpdateStatus(fmt.Sprintf("Cart processed: %d items, subtotal: $%.2f", cartCount, cartSubTotal))

		// * Create cart items from the summary items
		cartItems := make([]CartItem, 0)
		for _, item := range gjson.Get(string(body), "summaryItems").Array() {
			// * Format price with dollar sign
			price := item.Get("priceDetails.totalCustomerPrice").Float()
			totalPrice := fmt.Sprintf("$%.2f", price)

			// if price == 0 {
			// 	continue
			// }

			cartItems = append(cartItems, CartItem{
				Sku:                item.Get("skuId").String(),
				TotalCustomerPrice: totalPrice,
				IsDigital:          false,
			})
		}

		// * Create the checkout request
		checkoutReq := CartCheckoutRequest{
			ProductTotal:                   productTotalFormatted,
			SubTotal:                       cartSubTotal,
			HasSubscription:                false,
			IsAccountCreationRequired:      false,
			IsReadyForCheckoutWithDefaults: false,
			CartItems:                      cartItems,
		}

		data, err := json.Marshal(checkoutReq)
		if err != nil {
			return "", "", "", err
		}

		c.UpdateStatus("Item added to cart")

		return "", "", string(data), errors.New("no queue found")
	}

	c.UpdateStatus("Item is in queue")
	return a2cTransactionCode, a2cTransactionId, "", nil
}

// * Exit the queue by sending the a2cTransactionCode and a2cTransactionId in the 2nd add to cart request, and return the login URL from the response json
func (c *CheckoutClient) exitQueue(a2cTransactionCode, a2cTransactionId string) (string, string, error) {
	c.UpdateStatus("Exiting queue")

	resp, err := c.requestAddToCartQueue(a2cTransactionCode, a2cTransactionId)
	if err != nil {
		return "", "", err
	}

	a2cTransactionRefId := resp.Header.Get("a2ctransactionreferenceid")
	if a2cTransactionRefId == "" {
		return "", "", errors.New("failed to get a2c")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	loginUrl := gjson.Get(string(body), "redirectUrl").String()

	if loginUrl == "" {
		return "", "", errors.New("stuck in queue")
	}

	tokenSpl := strings.Split(loginUrl, "?token=")
	if len(tokenSpl) < 2 {
		return "", "", errors.New("stuck in queue")
	}

	token := tokenSpl[1]

	return token, a2cTransactionRefId, nil
}
func (c *CheckoutClient) getFinalAddToCart(a2cTransactionId string) (string, error) {
	resp, err := c.requestAddToCartFinal(a2cTransactionId)
	if err != nil {
		return "", err
	}

	c.UpdateStatus("Item added to cart")

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 412 {
		return "", errors.New("failed to add item to cart")
	}

	// Parse the cart response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var cartResponse CartResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&cartResponse)
	if err != nil {
		return "", err
	}

	// Store the cart response in the client
	c.UpdateStatus(fmt.Sprintf("Cart processed: %d items, subtotal: $%.2f", cartResponse.CartCount, cartResponse.CartSubTotal))

	productTotalStr := fmt.Sprintf("$%.2f", cartResponse.CartSubTotal)
	productTotalFormatted := strings.Replace(productTotalStr, ".", ",", 1)

	// Create cart items from the summary items
	cartItems := make([]CartItem, 0)
	for _, item := range cartResponse.SummaryItems {
		// Format price with dollar sign
		totalPrice := fmt.Sprintf("$%.2f", item.PriceDetails.TotalCustomerPrice)

		cartItems = append(cartItems, CartItem{
			Sku:                item.SkuId,
			TotalCustomerPrice: totalPrice,
			IsDigital:          false,
		})
	}

	// Create the checkout request
	checkoutReq := CartCheckoutRequest{
		ProductTotal:                   productTotalFormatted,
		SubTotal:                       cartResponse.CartSubTotal,
		HasSubscription:                false,
		IsAccountCreationRequired:      false,
		IsReadyForCheckoutWithDefaults: false,
		CartItems:                      cartItems,
	}

	data, err := json.Marshal(checkoutReq)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *CheckoutClient) carted(resp string) (bool, error) {
	itemsArr := gjson.Get(resp, "summaryItems").Array()

	if len(itemsArr) == 0 {
		return false, errors.New("failed to get items from cart")
	}

	for _, item := range itemsArr {
		if strings.EqualFold(item.Get("skuId").String(), c.Opts.SkuId) {
			return true, nil
		}
	}

	return false, nil
}
