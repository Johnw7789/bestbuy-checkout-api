package checkout

const (
	BestbuyAddToCartUrl    = "https://www.bestbuy.com/cart/api/v1/addToCart"
	BestbuyCartUrl         = "https://www.bestbuy.com/cart"
	BestbuyGoToCheckoutUrl = "https://www.bestbuy.com/cart/checkout"
	BestbuyCheckoutUrl     = "https://www.bestbuy.com/checkout/r/fast-track"

	BestbuyStoreLocatorUrl = "https://www.bestbuy.com/location/v1/US/store/zipcode/%s?types=Store&storeTypes=BigBox&status=Open"

	// Additional URLs for checkout steps
	BestbuyFulfillmentUrl         = "https://www.bestbuy.com/checkout/v3/%s"
	BestbuyPrelookupEndpoint      = "https://www.bestbuy.com/payment/api/v3/payment/%s/threeDSecure/preLookup"
	BestbuyPaymentKeyEndpoint     = "https://www.bestbuy.com/api/csiservice/v2/key/tas"
	BestbuyPaymentEndpoint        = "https://www.bestbuy.com/payment/api/v3/payment/%s/creditCard"
	BestbuyRefreshPaymentEndpoint = "https://www.bestbuy.com/checkout/v3/%s/paymentMethods/refreshPayments"
	Bestbuy3DSEndpoint            = "https://www.bestbuy.com/checkout/v3/paySecure/submitCardAuthentication"
	BestbuyPlaceOrderEndpoint     = "https://www.bestbuy.com/checkout/v3/%s"
	CardinalCommerceEndpoint      = "https://centinelapi.cardinalcommerce.com/V2/Cruise/StepUp"

	// URLs for checkout pages
	BestbuyShippingEndpoint    = "https://www.bestbuy.com/checkout/r/fulfillment"
	BestbuyPaymentPageEndpoint = "https://www.bestbuy.com/checkout/r/payment"
)
