package checkout

import (
	"github.com/Johnw7789/bestbuy-checkout/akamai"
	tls "github.com/bogdanfinn/tls-client"
)

type CheckoutClient struct {
	HttpClient     tls.HttpClient
	UpdateStatus   func(status string)
	Opts           CheckoutOpts
	AkamaiAdapter  *akamai.AkamaiAdapter

	// State fields for checkout process
	UserOrderId string
	OrderId     string
	PaymentId   string
	ItemIds     []string
	ThreeDsId   string
}

type CheckoutOpts struct {
	// Existing fields
	SkuId         string
	CVV           string
	Proxy         string
	UserAgent     string
	UserAgentHint string
	AkamaiApiKey  string

	// Profile fields
	Email       string
	PhoneNumber string

	// Shipping address fields
	ShippingFirstName string
	ShippingLastName  string
	ShippingAddress1  string
	ShippingAddress2  string
	ShippingCity      string
	ShippingStateCode string
	ShippingZipCode   string

	// Billing address fields
	BillingFirstName string
	BillingLastName  string
	BillingAddress1  string
	BillingAddress2  string
	BillingCity      string
	BillingStateCode string
	BillingZipCode   string

	// Payment fields
	CardNumber   string
	CardExpMonth string
	CardExpYear  string

	// Optional location ID for in-store pickup
	LocationID string
}

// ShipOrPickup represents the fulfillment type for an item
type FulfillmentPayload struct {
	User  User              `json:"user"`
	Items []ItemFulfillment `json:"items"`
}

type ItemFulfillment struct {
	Id          string      `json:"id"`
	Type        string      `json:"type"`
	Fulfillment Fulfillment `json:"fulfillment"`
}

type Fulfillment struct {
	Address ShipAddress `json:"address,omitempty"`
	Type    string      `json:"type"`
}

type Selectedfulfillment struct {
	Shipping      Shipping      `json:"shipping,omitempty"`
	InStorePickup InStorePickup `json:"inStorePickup,omitempty"`
}

type Shipping struct {
	Address ShipAddress `json:"address,omitempty"`
}

type InStorePickup struct {
	PickupStoreID         string `json:"pickupStoreId"`
	DisplayDateType       string `json:"displayDateType"`
	IsAvailableAtLocation bool   `json:"isAvailableAtLocation"`
	IsSTSAvailable        bool   `json:"isSTSAvailable"`
}

type ShipAddress struct {
	Addressline1         string `json:"addressLine1"`
	Addressline2         string `json:"addressLine2"`
	City                 string `json:"city"`
	Country              string `json:"country"`
	Firstname            string `json:"firstName"`
	Lastname             string `json:"lastName"`
	Middleinitial        string `json:"middleInitial"`
	Phonenumber          string `json:"phoneNumber"`
	Postalcode           string `json:"postalCode"`
	SaveAddressAsDefault bool   `json:"saveAddressAsDefault"`
	Savedtoprofile       bool   `json:"saveToProfile"`
	State                string `json:"state"`
	Type                 string `json:"type"`
	Useroverridden       bool   `json:"userOverridden"`
	Useforbilling        bool   `json:"useForBilling"`
	Classifer            string `json:"classifier"`
}

type ContactInfo struct {
	Phonenumber     string     `json:"phoneNumber"`
	Smsnotifynumber string     `json:"smsNotifyNumber"`
	Smsoptin        bool       `json:"smsOptIn"`
	Emailaddress    string     `json:"emailAddress"`
	Items           []ShipItem `json:"items"`
}

type ShipItem struct {
	Type                string              `json:"type"`
	ID                  string              `json:"id"`
	Selectedfulfillment Selectedfulfillment `json:"selectedFulfillment"`
	Giftmessageselected bool                `json:"giftMessageSelected"`
}

type PaymentRequest struct {
	Billingaddress Billingaddress `json:"billingAddress"`
	Creditcard     Creditcard     `json:"creditCard"`
}

type Billingaddress struct {
	Addressline1    string `json:"addressLine1"`
	Addressline2    string `json:"addressLine2"`
	City            string `json:"city"`
	Country         string `json:"country"`
	Dayphone        string `json:"dayPhone"`
	Emailaddress    string `json:"emailAddress"`
	Firstname       string `json:"firstName"`
	IsSimpleAddress bool   `json:"isSimpleAddress"`
	Lastname        string `json:"lastName"`
	Postalcode      string `json:"postalCode"`
	Standardized    bool   `json:"standardized"`
	State           string `json:"state"`
	Useroverridden  bool   `json:"userOverridden"`
}

type Creditcard struct {
	AppliedPoints       string `json:"appliedPoints"`
	Binnumber           string `json:"binNumber"`
	CreditCardProfileId string `json:"creditCardProfileId"`
	Cvv                 string `json:"cvv"`
	Default             bool   `json:"default"`
	Expmonth            string `json:"expMonth"`
	Expyear             string `json:"expYear"`
	Number              string `json:"number"`
	Orderid             string `json:"orderId"`
	PurchaseOrderNumber string `json:"purchaseOrderNumber"`
	Savetoprofile       bool   `json:"saveToProfile"`
	Type                string `json:"type"`
	Virtualcard         bool   `json:"virtualCard"`
}

type PrelookupRequest struct {
	BinNumber           string      `json:"binNumber"`
	Browserinfo         Browserinfo `json:"browserInfo"`
	Orderid             string      `json:"orderId"`
	PayComponentVersion string      `json:"payComponentVersion"`
	PaymentId           string      `json:"paymentId"`
	PayVersion          string      `json:"payVersion"`
}

type Browserinfo struct {
	Colordepth  string `json:"colorDepth"`
	Height      string `json:"height"`
	Javaenabled bool   `json:"javaEnabled"`
	Language    string `json:"language"`
	Timezone    string `json:"timeZone"`
	Useragent   string `json:"userAgent"`
	Width       string `json:"width"`
}

type PrelookupResponse struct {
	ThreeDSReferenceId      string `json:"threeDSReferenceId"`
	DeviceDataCollectionJwt string `json:"deviceDataCollectionJwt"`
	DeviceDataCollectionUrl string `json:"deviceDataCollectionUrl"`
	PaymentId               string `json:"paymentId"`
}

type PlaceOrderRequest struct {
	User     User     `json:"user"`
	Metadata Metadata `json:"metadata"`
}

type User struct{}

type Metadata struct {
	Userdevicedetails Userdevicedetails `json:"userDeviceDetails"`
	Browserinfo       Browserinfo       `json:"browserInfo"`
	Buylegalmessaging []string          `json:"buyLegalMessaging"`
}

type Userdevicedetails struct {
	ForterToken string `json:"forterToken"`
}

type Threedsecurestatus struct {
	Threedsreferenceid string `json:"threeDSReferenceId"`
}

type ErrorResponse struct {
	Errors []struct {
		Errorcode string `json:"errorCode"`
		Message   string `json:"message"`
	} `json:"errors"`
}

// Cart response from Add To Cart request
type CartResponse struct {
	CartCount    int           `json:"cartCount"`
	CartSubTotal float64       `json:"cartSubTotal"`
	SummaryItems []SummaryItem `json:"summaryItems"`
	Category     string        `json:"category"`
}

type SummaryItem struct {
	SkuId                string       `json:"skuId"`
	ZipCode              string       `json:"zipCode"`
	IsICR                bool         `json:"isICR"`
	IsHaccs              bool         `json:"isHaccs"`
	IsPaidMemberDiscount bool         `json:"isPaidMemberDiscount"`
	LineId               string       `json:"lineId"`
	Quantity             int          `json:"quantity"`
	Price                float64      `json:"price"`
	PriceDetails         PriceDetails `json:"priceDetails"`
	RankingId            string       `json:"rankingId"`
}

type PriceDetails struct {
	PriceEventType     string  `json:"priceEventType"`
	TotalCurrentPrice  float64 `json:"totalCurrentPrice"`
	TotalCustomerPrice float64 `json:"totalCustomerPrice"`
	TotalRegularPrice  float64 `json:"totalRegularPrice"`
	UnitCurrentPrice   float64 `json:"unitCurrentPrice"`
	UnitRegularPrice   float64 `json:"unitRegularPrice"`
	GspUnitPrice       float64 `json:"gspUnitPrice"`
	ListPrice          float64 `json:"listPrice"`
	COPUnitPrice       float64 `json:"COPUnitPrice"`
}

// Struct for the cart checkout request
type CartCheckoutRequest struct {
	ProductTotal                   string     `json:"productTotal"`
	SubTotal                       float64    `json:"subTotal"`
	HasSubscription                bool       `json:"hasSubscription"`
	IsAccountCreationRequired      bool       `json:"isAccountCreationRequired"`
	IsReadyForCheckoutWithDefaults bool       `json:"isReadyForCheckoutWithDefaults"`
	CartItems                      []CartItem `json:"cartItems"`
}

type CartItem struct {
	Sku                string `json:"sku"`
	TotalCustomerPrice string `json:"totalCustomerPrice"`
	IsDigital          bool   `json:"isDigital"`
}

type Order struct {
	OrderId    string  `json:"orderId"`
	TotalPrice float64 `json:"totalPrice"`
	PickupDate string  `json:"pickupDate"`
}

type CardAuthentication struct {
	TransactionId string `json:"TransactionId"`
	PaRes string `json:"PaRes"`
	OrderId 	string `json:"orderId"`
}