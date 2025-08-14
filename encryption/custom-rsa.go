package encryption

import (
	"encoding/binary"
	"errors"
	"hash"
	"io"
)

// generates a mask using the MGF1 algorithm
func mgf1(seed []byte, length int, hashFunc func() hash.Hash) []byte {
	t := []byte{}

	for counter := uint32(0); len(t) < length; counter++ {
		c := make([]byte, 4)
		binary.BigEndian.PutUint32(c, counter)

		h := hashFunc()

		h.Write(seed)
		h.Write(c)

		t = append(t, h.Sum(nil)...)
	}

	return t[:length]
}

// pad with separate hash functions for md padding and MGF1
func oaepPad(rand io.Reader, hashFunc func() hash.Hash, mgfHashFunc func() hash.Hash, label, msg []byte, k int) ([]byte, error) {
	h := hashFunc()
	h.Write(label)
	lHash := h.Sum(nil)
	hLen := h.Size()

	if len(msg) > k-2*hLen-2 {
		return nil, errors.New("message too long for RSA key size")
	}

	seed := make([]byte, hLen)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, err
	}

	psLen := k - len(msg) - 2*hLen - 2

	db := make([]byte, 0, k-hLen-1)
	db = append(db, lHash...)
	db = append(db, make([]byte, psLen)...) // PS is zero bytes
	db = append(db, 0x01)
	db = append(db, msg...)

	mgfDB := mgf1(seed, len(db), mgfHashFunc)
	maskedDB := make([]byte, len(db))
	for i := range db {
		maskedDB[i] = db[i] ^ mgfDB[i]
	}

	mgfMaskedDB := mgf1(maskedDB, hLen, mgfHashFunc)
	maskedSeed := make([]byte, hLen)
	for i := range seed {
		maskedSeed[i] = seed[i] ^ mgfMaskedDB[i]
	}

	em := make([]byte, 0, k)
	em = append(em, 0x00)
	em = append(em, maskedSeed...)
	em = append(em, maskedDB...)

	return em, nil
}
