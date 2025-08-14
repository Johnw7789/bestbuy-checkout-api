package checkout

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// * Decode the queue time from the a2cTransactionCode
func (c *CheckoutClient) decodeQueue(a2cTransactionCode string) (int, error) {
	if len(a2cTransactionCode) < 7 {
		return 0, errors.New("invalid a2cTransactionCode")
	}

	prefixHex := a2cTransactionCode[:3]
	iterCount, err := strconv.ParseInt(prefixHex, 16, 64)
	if err != nil {
		return 0, errors.New("failed to parse iterCount")
	}

	codeSubstring := a2cTransactionCode[3:7]
	targetHash := a2cTransactionCode[7:]

	// * Try all potential values of time (in seconds) in the range of 0 to 899, the max potential queue time that bb seems to have
	for potentialTime := 0; potentialTime < 900; potentialTime++ {
		// * Create the hash
		toHash := fmt.Sprintf("%d%s%s", potentialTime, codeSubstring, c.Opts.SkuId)
		hashedBytes := []byte(toHash)

		for i := int64(0); i < iterCount; i++ {
			sum := sha256.Sum256(hashedBytes)
			hashedBytes = sum[:]
		}

		// * Convert final hash to hex and compare
		if hex.EncodeToString(hashedBytes) == targetHash {
			return potentialTime, nil // * We found the matching queue time
		}
	}

	return 0, errors.New("no matching queue time found")
}
