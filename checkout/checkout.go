package checkout

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Johnw7789/bestbuy-checkout/encryption"
	"github.com/Johnw7789/bestbuy-checkout/login"
)

// * Handle the cart process. If there is a queue, it will handle the queue and return the tokenId and a2cTransactionRefId to be used in Login()/Checkout()
func (c *CheckoutClient) CartItem() (string, string, string, error) {
	c.UpdateStatus("Starting cart process")

	homePage, err := c.reqHomePage()
	if err != nil {
		return "", "", "", err
	}

	// * Check if akamai is blocking us with sbsd site script
	doSbsd := false
	if !strings.Contains(homePage, "skuId") {
		doSbsd = true
	}

	// * Generate abck cookie
	err = c.handleAkamai(homePage, doSbsd)
	if err != nil {
		return "", "", "", err
	}

	// * 1st add to cart request
	a2cTransactionCode, a2cTransactionId, checkoutData, err := c.tryCart()
	if err != nil {
		// * Will 99% be an error saying no queue. So return checkoutData and proceed to login/checkout
		return "", "", checkoutData, err
	}

	// * Deciper the queue time from the a2cTransactionCode response header
	queueTime, err := c.decodeQueue(a2cTransactionCode)
	if err != nil {
		return "", "", "", errors.New("failed to decode queue")
	}

	c.UpdateStatus("Queue time: " + strconv.Itoa(queueTime) + " seconds")

	// * Add the necessary 1 second buffer to ensure we don't try to "exit" the queue too early and get a new queue time
	time.Sleep(time.Duration(queueTime+1) * time.Second)

	// * 2/3 add to cart request
	tokenId, a2cTransactionRefId, err := c.exitQueue(a2cTransactionCode, a2cTransactionId)
	if err != nil {
		return "", "", "", err
	}

	c.UpdateStatus("Queue exited, logging in")

	return tokenId, a2cTransactionRefId, checkoutData, nil
}

// * Checkout takes a login client that has already been authenticated and handles the full checkout process after login and after queue (if there was one)
func (c *CheckoutClient) Checkout(loginClient *login.LoginClient, redirectUrl, a2cId, checkoutData string) (Order, error) {
	c.UpdateStatus("Starting checkout process")

	// * The redirect uri returned from the final login response
	prodPage, err := c.reqProdPageRedirect(redirectUrl)
	if err != nil {
		return Order{}, err
	}

	// * Check if we need to handle akamai sbsd by seeing if we are blocked by the sbsd script when access the prod page
	doSbsd := false
	if !strings.Contains(prodPage, "skuId") {
		doSbsd = true
	}

	// todo: evaluate whether or not we need to handle akamai for a 3rd time during this process - could be overkill??
	c.UpdateStatus("Handling Akamai")
	err = c.handleAkamai(prodPage, doSbsd)
	if err != nil {
		return Order{}, err
	}

	if a2cId != "" {
		c.UpdateStatus("Adding to cart")
		// * Send the 3rd and final add to cart request to complete the queue process (needs to be done after login)
		checkoutData, err = c.getFinalAddToCart(a2cId)
		if err != nil {
			return Order{}, err
		}
	}

	c.UpdateStatus("Submitting cart")
	// * Send the initiate checkout request
	checkoutResp, err := c.reqCartCheckout(checkoutData)
	if err != nil {
		return Order{}, err
	}

	if checkoutResp.StatusCode != 200 {
		return Order{}, fmt.Errorf("failed to process cart checkout, status: %d", checkoutResp.StatusCode)
	}

	c.UpdateStatus("Getting checkout data")
	// * Get order details from the checkout page
	shippingAvailable, err := c.getOrderDetails()
	if err != nil {
		return Order{}, err
	}

	// todo: use the reqClosestStores function to loop through all stores within designated radius
	// for now though, if an item is pickup only, it seems that bestbuy should set cookies that force your closest store/closest with item in stock 
	// so we can add this later, but also prioritize shipping if the option is available
	if shippingAvailable {
		c.UpdateStatus("Submitting shipping info")
		err = c.getSetFufillment()
		if err != nil {
			return Order{}, err
		}
	}

	// Set contact and shipping information
	// err = c.getSetContactInfo()
	// if err != nil {
	// 	return Order{}, err
	// }

	c.UpdateStatus("Getting payment data")
	// Perform payment prelookup to get 3DS reference ID
	err = c.getPaymentPrelookup()
	if err != nil {
		return Order{}, err
	}

	pubKey, _, err := c.getPublicKey(BestbuyPaymentKeyEndpoint)
	if err != nil {
		return Order{}, err
	}

	encCard, err := encryption.EncryptCardNumber(c.Opts.CardNumber, pubKey)
	if err != nil {
		return Order{}, err
	}

	c.UpdateStatus("Submitting payment")
	// * Set payment information
	maskedNumber, err := c.getSetPaymentInfo(encCard + ":3:735818052:" + c.Opts.CardNumber)
	if err != nil {
		return Order{}, err
	}

	c.UpdateStatus("Refreshing payment")
	// * Refresh payment options
	err = c.getRefreshPayment()
	if err != nil {
		return Order{}, err
	}

	c.UpdateStatus("Submitting payment (2)")
	// * Set payment information again (might not be needed)
	_, err = c.getSetPaymentInfo(maskedNumber) // todo: check if needed
	if err != nil {
		return Order{}, err
	}

	c.UpdateStatus("Submitting order")
	// * Handle final payment/3ds if needed and then placing the order
	order, err := c.getPlaceOrder()
	if err != nil {
		return Order{}, err
	}

	c.UpdateStatus("Order placed successfully!")

	return order, nil
}
