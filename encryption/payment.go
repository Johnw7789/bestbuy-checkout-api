package encryption

// * Encrypt cc num using Bestbuy's prefix "00960001"
func EncryptCardNumber(cardNumber, pubKey string) (string, error) {
	encCard, err := Encrypt("00960001"+cardNumber, pubKey)
	if err != nil {
		return "", err
	}

	return encCard, nil
}
