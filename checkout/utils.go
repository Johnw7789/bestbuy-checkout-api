package checkout

import (
	"fmt"
	"time"
)

// * Get the card type based on the first digit of the card number
func getCardType(cardNumber string) string {
	if len(cardNumber) == 0 {
		return "VISA"
	}

	firstDigit := cardNumber[0]
	switch firstDigit {
	case '4':
		return "VISA"
	case '5':
		return "MASTERCARD"
	case '3':
		return "AMEX"
	case '6':
		return "DISCOVER"
	default:
		return "VISA"
	}
}

func generateUUID() string {
	// Create a UUID based on the current time
	return fmt.Sprintf("%s-f3e3-11ef-bcdd-3bd1d3b23d12", fmt.Sprintf("%x", time.Now().UnixNano())[0:8])
}
