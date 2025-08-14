package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"strings"
)

func Encrypt(s string, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	var pubkey *rsa.PublicKey
	pubkey, _ = parsedKey.(*rsa.PublicKey)

	// encrypt using sha1 for hash and for MGF1
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(s), nil)
	if err != nil {
		return "", err
	}

	// return in bestbuy format
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func EncryptJoin(s string, publicKey string, keyId string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	var pubkey *rsa.PublicKey
	pubkey, _ = parsedKey.(*rsa.PublicKey)

	// encrypt using sha1 for hash and for MGF1
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(s), nil)
	if err != nil {
		return "", err
	}

	// return in bestbuy format
	return strings.Join([]string{"1", keyId, base64.StdEncoding.EncodeToString(ciphertext)}, ":"), nil
}

func EncryptSmall(s string, publicKey string, keyId string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	var pubkey *rsa.PublicKey
	pubkey, _ = parsedKey.(*rsa.PublicKey)

	// calc modulus size in bytes
	k := (pubkey.N.BitLen() + 7) / 8

	// add bestbuy specific padding
	msg := "4PADDING" + s

	// apply OAEP padding - sha256 for hash, sha1 for MGF1
	em, err := oaepPad(rand.Reader, sha256.New, sha1.New, []byte{}, []byte(msg), k)
	if err != nil {
		return "", err
	}

	// encrypt
	m := new(big.Int).SetBytes(em)
	c := new(big.Int).Exp(m, big.NewInt(int64(pubkey.E)), pubkey.N)
	cBytes := c.Bytes()

	// make sure ciphertext is k bytes by padding with leading zeros if necessary
	if len(cBytes) < k {
		padded := make([]byte, k)
		copy(padded[k-len(cBytes):], cBytes)
		cBytes = padded
	}

	// return in bestbuy format
	return strings.Join([]string{"1", keyId, base64.StdEncoding.EncodeToString(cBytes), "1"}, ":"), nil
}
