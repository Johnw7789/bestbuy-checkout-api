package login

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strconv"
	"time"
)

func getOTP(secret2fa string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(secret2fa)
	if err != nil {
		return "", err
	}

	interval := time.Now().Unix() / 30

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	hash := hmac.New(sha1.New, key)
	hash.Write(bs)

	h := hash.Sum(nil)
	o := (h[19] & 15)

	var header uint32
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return "", err
	}

	h12 := (int(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(int(h12))

	if len(otp) == 6 {
		return otp, nil
	}

	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}

	return otp, nil
}
